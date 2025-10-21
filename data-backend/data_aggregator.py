from os import listdir
from os.path import isfile, join
import polars as pl
from google.cloud import storage
import io

prices_path = "./prices"
volumes_path = "./volumes"
client = storage.Client()
bucket = client.bucket("house-prices-volume-dataset")
is_prod = True

volumes_aliases = {
    0 : "Range <50 m²",
    50 : "Range 50 - 85 m²",
    85 : "Range 85 - 115 m²",
    115 : "Range 115 - 145 m²",
    145 : "Range >145 m²",
}

def _get_file_list_from_bucket(prefix:str) -> list[str]:

    # Use prefix to simulate directories, e.g., "data/"
    blobs = bucket.list_blobs(prefix=prefix)

    # Filter out directories (blobs ending with "/") and sort by name
    files = [blob.name for blob in blobs if not blob.name.endswith("/")]
    files.sort()

    print(files)

    return files


def _get_file_list(path: str):
    files = [f for f in listdir(path) if isfile(join(path, f))]
    files.sort()

    return files

prices_file_list = _get_file_list_from_bucket("prices/") if is_prod else _get_file_list(prices_path)
volumes_file_list = _get_file_list_from_bucket("volumes/") if is_prod else _get_file_list(volumes_path)

def get_price_starting_year():
    oldest_file = prices_file_list[0]
    year = oldest_file[:4]
    semester = oldest_file[5:7]
    return semester + "/" + year

def get_price_current_year():
    last_file = prices_file_list[-1]
    year = last_file[:4]
    semester = last_file[5:7]
    return semester + "/" + year

def get_volume_starting_year():
    if is_prod:
        return int(volumes_file_list[0].split("/")[1][:4])
    else:
        return int(volumes_file_list[0][:4])

def get_volume_current_year():
    if is_prod:
        return int(volumes_file_list[-1].split("/")[1][:4])
    else:
        return int(volumes_file_list[-1][:4])


def _get_year(s: str) -> str:
    if is_prod:
        return s.split("/")[1][:4]
    else:
        return s[:4]

def _get_year_semester(s: str) -> str:
    if is_prod:
        return s.split("/")[1][:7]
    else:
        return s[:7]


def _semester_converter(s: str) -> str:
    return (s[5:7] if s[5:7] == "01" else "06") + "/" + s[:4]


def get_number_of_price_files() -> int:
    return len([f for f in prices_file_list if "VALORI" in f])

def get_price_df():
    # Prices section
    files = prices_file_list
    prices_files = [f for f in files if "VALORI" in f]

    prices = pl.DataFrame()
    cols = [
        "Comune_amm",
        "Comune_descrizione",
        "Zona",
        "Descr_Tipologia",
        "Cod_Tip",
        "Stato",
        "Compr_min",
        "Compr_max"
    ]

    for file in prices_files:
        if is_prod:
            blob = bucket.blob(file)
            data = blob.download_as_bytes()
            prices_vals = pl.read_csv(io.BytesIO(data), separator=";", columns=cols, skip_rows=1)
        else:
            prices_vals = pl.read_csv(join(prices_path, file),
                separator=";", columns=cols, skip_rows=1
            )

        prices_vals =prices_vals.filter(
            pl.col("Cod_Tip").is_in([1, 13, 14, 15, 19, 20, 21, 22])
        )

        # Rename columns
        prices_vals = prices_vals.rename({
            "Comune_amm": "COD_COMUNE",
            "Comune_descrizione": "DES_COMUNE",
            "Zona": "COD_ZONA",
            "Descr_Tipologia": "DES_TIPOLOGIA",
            "Stato": "DES_STATO",
            "Compr_min": "COSTO_MIN",
            "Compr_max": "COSTO_MAX"
        })

        prices_vals = prices_vals.with_columns(
            pl.lit(_semester_converter(_get_year_semester(file))).alias("MESE_ANNO")
        )

        if prices.is_empty():
            prices = prices_vals
        else:
            prices = pl.concat([prices_vals, prices], how="vertical")


    current_zone_file = [f for f in files if "ZONE" in f][-1]

    path = ""
    prices_infos = None
    if is_prod:
        path = current_zone_file
        blob = bucket.blob(path)
        data = blob.download_as_bytes()
        prices_infos = pl.read_csv(io.BytesIO(data),
            separator=";",
            columns=[
                "Comune_amm",
                "Zona",
                "Zona_Descr",
            ],
            skip_rows=1
        )
    else:
        path = join(prices_path, current_zone_file)
        prices_infos = pl.read_csv(path,
            separator=";",
            columns=[
                "Comune_amm",
                "Zona",
                "Zona_Descr",
            ],
            skip_rows=1
        )


    # Rename columns
    prices_infos = prices_infos.rename({
        "Comune_amm": "COD_COMUNE",
        "Zona": "COD_ZONA",
        "Zona_Descr": "DES_ZONA",
    })

    prices_infos = prices_infos.with_columns(
        pl.col("DES_ZONA").str.strip_chars("'").alias("DES_ZONA")
    )

    prices = prices_infos.join(prices, on=["COD_COMUNE", "COD_ZONA"], how="left").drop_nulls()
    return prices.select(
        "MESE_ANNO",
        "DES_COMUNE",
        "DES_ZONA",
        "DES_TIPOLOGIA",
        "DES_STATO",
        "COSTO_MIN",
        "COSTO_MAX"
    )

def get_volume_df():
    # Volume section
    files = volumes_file_list
    volumes_files = [f for f in files if "RES" in f]

    volumes = pl.DataFrame()
    for file in volumes_files:
        if is_prod:
            blob = bucket.blob(file)
            data = blob.download_as_bytes()
            volumes_vals = pl.read_csv(io.BytesIO(data), separator=";")
        else:
            volumes_vals = pl.read_csv(join(volumes_path, file),
                separator=";"
            )

        volumes_vals = volumes_vals.rename(
            {col: col.upper() for col in volumes_vals.columns}
        )

        # Getting detailed information
        volumes_vals = volumes_vals.select(
            pl.col(_get_year(file) + "_CODCOM").alias("COD_COMUNE"),
            pl.col("NTN " + _get_year(file) + " FINO A 50 MQ").alias(volumes_aliases[0]),
            pl.col("NTN " + _get_year(file) + " 50 -| 85 MQ").alias(volumes_aliases[50]),
            pl.col("NTN " + _get_year(file) + " 85 -| 115 MQ").alias(volumes_aliases[85]),
            pl.col("NTN " + _get_year(file) + " 115 -| 145 MQ").alias(volumes_aliases[115]),
            pl.col("NTN " + _get_year(file) + " OLTRE 145 MQ").alias(volumes_aliases[145]),
            pl.col("NTN_" + _get_year(file)).alias("SOMMA_NTN")
        )


        volumes_vals = volumes_vals.with_columns(
            pl.lit(_get_year(file)).alias("ANNO")
        )

        if volumes.is_empty():
            volumes = volumes_vals
        else:
            volumes = pl.concat([volumes_vals, volumes], how="vertical")


    current_com_file = [f for f in files if "COM" in f][-1]

    path = ""
    volumes_detail = None
    if is_prod:
        path = current_com_file
        blob = bucket.blob(path)
        data = blob.download_as_bytes()
        volumes_detail = pl.read_csv(io.BytesIO(data),
            separator=";",
        )
    else:
        path = join(volumes_path, current_com_file)
        volumes_detail = pl.read_csv(path,
            separator=";"
        )


    volumes_detail = volumes_detail.rename(
        {col: col.upper() for col in volumes_detail.columns}
    )

    volumes_detail = volumes_detail.select(
        pl.col(_get_year(current_com_file) + "_CODCOM").alias("COD_COMUNE"),
        pl.col("COMUNE").alias("DES_COMUNE"),
    )

    volumes = volumes_detail.join(volumes, on="COD_COMUNE", how='left').drop_nulls()

    return volumes.select(
        "ANNO",
        "DES_COMUNE",
        "Range <50 m²",
        "Range 50 - 85 m²",
        "Range 85 - 115 m²",
        "Range 115 - 145 m²",
        "Range >145 m²",
        "SOMMA_NTN"
    )

def get_volume_market_size_df():
    # Getting the informations about the size of the market
    # Currently this information is only in the 2024 file (volumes/2024_VALORI_RES.csv)
    # So i took the last one
    market_size_file = [f for f in volumes_file_list if "COM" in f][-1]

    path = ""
    volumes_market_size = None
    if is_prod:
        path = market_size_file
        blob = bucket.blob(path)
        data = blob.download_as_bytes()
        volumes_market_size = pl.read_csv(io.BytesIO(data),
            separator=";",
        )
    else:
        path = join(volumes_path, market_size_file)
        volumes_market_size = pl.read_csv(path,
            separator=";"
        )

    volumes_market_size = volumes_market_size.rename(
        {col: col.upper() for col in volumes_market_size.columns}
    )

    # Getting the informations about the size of the market
    volumes_market_size = volumes_market_size.select(
        pl.col(_get_year(market_size_file) + "_CODCOM").alias("COD_COMUNE"),
        pl.col("COMUNE").alias("DES_COMUNE"),
        pl.col("TAGLIA MERCATO").alias("TAGLIA_MERCATO")
    ).drop_nulls()


    return volumes_market_size.select(
        "DES_COMUNE",
        "TAGLIA_MERCATO"
    )

def get_municipality_info():
    return get_price_df().select(
        "DES_COMUNE",
        "DES_ZONA",
        "DES_TIPOLOGIA",
        "DES_STATO"
    ).unique(keep='first').sort("DES_COMUNE")
