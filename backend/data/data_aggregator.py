from os import listdir
from os.path import isfile, join
import polars as pl

prices_path = "./prices"
volumes_path = "./volumes"

def _get_year(s: str) -> str:
    return s[:4]

def _get_year_semester(s: str) -> str:
    return s[:7]

def _semester_converter(s: str) -> str:
    return (s[5:7] if s[5:7] == "01" else "06") + "/" + s[:4]

def get_prices_df():
    # Prices section
    prices_files = [f for f in listdir(prices_path) if isfile(join(prices_path, f))]
    prices_files.sort()

    prices = pl.DataFrame()
    for i, file in enumerate(prices_files[:-1]):
        next_file = prices_files[i + 1]

        # If year_semester is the same, join the dataframe together
        if _get_year_semester(file) == _get_year_semester(next_file):
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

            prices_infos = pl.read_csv(join(prices_path, next_file),
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

            prices_vals = prices_vals.join(prices_infos, on=["COD_COMUNE", "COD_ZONA"], how="left")

            prices_vals = prices_vals.with_columns(
                pl.lit(_semester_converter(_get_year_semester(file))).alias("MESE_ANNO")
            )

            if prices.is_empty():
                prices = prices_vals
            else:
                prices = pl.concat([prices, prices_vals], how="vertical")

    return prices

def get_volumes_df():
    # Volume section
    volumes_files = [f for f in listdir(volumes_path) if isfile(join(volumes_path, f))]
    volumes_files.sort()

    volumes = pl.DataFrame()

    for i, file in enumerate(volumes_files[:-1]):
        next_file = volumes_files[i + 1]

        if _get_year(file) == _get_year(next_file):


            volumes_market_size = pl.read_csv(join(volumes_path, file),
                separator=";"
            )
            volumes_market_size = volumes_market_size.rename(
                {col: col.upper() for col in volumes_market_size.columns}
            )

            # Getting the informations about the size of the market
            volumes_market_size = volumes_market_size.select(
                pl.col(_get_year(file) + "_CODCOM").alias("COD_COMUNE"),
                pl.col("COMUNE").alias("DES_COMUNE")
            )

            volumes_detail = pl.read_csv(join(volumes_path, next_file),
                separator=";"
            )
            volumes_detail = volumes_detail.rename(
                {col: col.upper() for col in volumes_detail.columns}
            )

            # Getting detailed information
            volumes_detail = volumes_detail.select(
                pl.col(_get_year(file) + "_CODCOM").alias("COD_COMUNE"),
                pl.col("NTN " + _get_year(file) + " FINO A 50 MQ").alias("<50_MQ"),
                pl.col("NTN " + _get_year(file) + " 50 -| 85 MQ").alias("<85_MQ"),
                pl.col("NTN " + _get_year(file) + " 85 -| 115 MQ").alias("<115_MQ"),
                pl.col("NTN " + _get_year(file) + " 115 -| 145 MQ").alias("<145_MQ"),
                pl.col("NTN " + _get_year(file) + " OLTRE 145 MQ").alias(">145_MQ"),
                pl.col("NTN_" + _get_year(file)).alias("SOMMA_COLUMI")
            )

            volumes_detail = volumes_detail.join(volumes_market_size, on="COD_COMUNE", how='inner')

            volumes_detail = volumes_detail.with_columns(
                pl.lit(_get_year(file)).alias("ANNO")
            )

            if volumes.is_empty():
                volumes = volumes_detail
            else:
                volumes = pl.concat([volumes, volumes_detail], how="vertical")

    return volumes

def get_volumes_market_size_df():
    # Getting the informations about the size of the market
    # Currently this information is only in the 2024 file (volumes/2024_VALORI_RES.csv)
    volumes_files = [f for f in listdir(volumes_path) if isfile(join(volumes_path, f)) and ("LISTA_COM.csv" in f)]
    volumes_files.sort()

    # Get the last LISTA_COM.csv file loaded in the folder
    last_csv_name = volumes_files[-1]
    year = _get_year(last_csv_name)

    volumes_market_size = pl.read_csv(join(volumes_path, last_csv_name),
        separator=";"
    )
    volumes_market_size = volumes_market_size.rename(
        {col: col.upper() for col in volumes_market_size.columns}
    )

    # Getting the informations about the size of the market
    volumes_market_size = volumes_market_size.select(
        pl.col(year + "_CODCOM").alias("COD_COMUNE"),
        pl.col("COMUNE").alias("DES_COMUNE"),
        pl.col("TAGLIA MERCATO").alias("TAGLIA_MERCATO")
    )

    return volumes_market_size
