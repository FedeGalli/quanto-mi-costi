from os import listdir
from os.path import isfile, join
import polars as pl

prices_path = "./prices"
volumes_path = "./volumes"
volumes_aliases = {
    0 : "Range <50 m²",
    50 : "Range 50 - 85 m²",
    85 : "Range 85 - 115 m²",
    115 : "Range 115 - 145 m²",
    145 : "Range >145 m²",
}

def _get_file_list(path: str):
    files = [f for f in listdir(path) if isfile(join(path, f))]
    files.sort()

    return files

def get_price_starting_year():
    oldest_file = _get_file_list(prices_path)[0]
    year = oldest_file[:4]
    semester = oldest_file[5:7]
    return semester + "/" + year

def get_price_current_year():
    last_file = _get_file_list(prices_path)[-1]
    year = last_file[:4]
    semester = last_file[5:7]
    return semester + "/" + year

def get_volume_starting_year():
    return int(_get_file_list(volumes_path)[0][:4])

def get_volume_current_year():
    return int(_get_file_list(volumes_path)[-1][:4])

def _get_year(s: str) -> str:
    return s[:4]

def _get_year_semester(s: str) -> str:
    return s[:7]

def _semester_converter(s: str) -> str:
    return (s[5:7] if s[5:7] == "01" else "06") + "/" + s[:4]

def get_number_of_price_files() -> int:
    return len([f for f in _get_file_list(prices_path) if "VALORI" in f])

def get_price_df():
    # Prices section
    files = _get_file_list(prices_path)
    prices_files = [f for f in files if "VALORI" in f]

    prices = pl.DataFrame()
    for file in prices_files:
        prices_vals = pl.read_csv(join(prices_path, file),
            separator=";",
            columns=[
                "Comune_amm",
                "Comune_descrizione",
                "Zona",
                "Descr_Tipologia",
                "Cod_Tip",
                "Stato",
                "Compr_min",
                "Compr_max"
            ],
            skip_rows=1
        ).filter(
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

    prices_infos = pl.read_csv(join(prices_path, current_zone_file),
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
    files = _get_file_list(volumes_path)
    volumes_files = [f for f in files if "RES" in f]

    volumes = pl.DataFrame()
    for file in volumes_files:
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


    volumes_detail = pl.read_csv(join(volumes_path, current_com_file),
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
    market_size_file = [f for f in _get_file_list(volumes_path) if "COM" in f][-1]

    volumes_market_size = pl.read_csv(join(volumes_path, market_size_file),
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
