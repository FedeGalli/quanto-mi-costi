import polars as pl
from data_aggregator import get_volume_market_size_df, get_municipality_info, get_volume_df, get_price_df, get_price_current_year, get_price_starting_year, get_volume_current_year, get_volume_starting_year

class DataManager:
    def __init__(self):
        self.price_start_year = get_price_starting_year()
        self.price_current_year = get_price_current_year()
        self.volume_start_year = get_volume_starting_year()
        self.volume_current_year = get_volume_current_year()
        self.volumes_data = get_volume_df()
        self.prices_data = get_price_df()
        self.market_size_data = get_volume_market_size_df()
        self.municipality_info = get_municipality_info()

    def _get_volume_filtered(self, com: str, year:str):
        return self.volumes_data.filter((pl.col("DES_COMUNE") == com) & (pl.col("ANNO") == year))

    def _get_prices_filtered(self, com: str, month_year:str):
        return self.prices_data.filter((pl.col("DES_COMUNE") == com) & (pl.col("MESE_ANNO") == month_year))

    def get_current_volume_mq(self, com: str, mq:int):
        # Retrieving columns name with number in the text
        allowed_mq_columns = [f for f in self.volumes_data.columns if "<" in f or ">" in f]
        mq_bounds = [int(f[1:-3]) for f in allowed_mq_columns]

        # Getting the right column base on the MQ passed in the function
        selected_mq_col = None
        for i, col in enumerate(allowed_mq_columns):
            if mq < mq_bounds[i] or i == len(allowed_mq_columns) - 1:
                selected_mq_col = col
                break

        return self._get_volume_filtered(com, str(self.volume_current_year)).select(selected_mq_col)

    def get_current_price(self, com:str, zone:str, type: str, state:str):
        max_and_min = self._get_prices_filtered(com, str(self.price_current_year)).filter(
            (pl.col("DES_ZONA") == zone) &
            (pl.col("DES_STATO") == state) &
            (pl.col("DES_TIPOLOGIA") == type)
        ).select("COSTO_MIN", "COSTO_MAX")
        return max_and_min.select("COSTO_MIN"), max_and_min.select("COSTO_MAX")

    def get_price_trend(self, com: str, zone:str, type: str, state:str):
        trend_df = self._get_prices_filtered(com, str(self.price_start_year))
        print(int(self.price_current_year[3:]))
        for i in range((int(self.price_current_year[3:]) - int(self.price_start_year[3:])) * 2):
            print("CIAO")
            print("0" + str(((i % 2) * 5) + 1) + "/" + str(int(self.price_start_year[3:]) + int(i / 2) + 1))
            trend_df = pl.concat(
                [
                    # Building the semester/year with variable i
                    self._get_prices_filtered(com, "0" + str(((i % 2) * 5) + 1) + "/" + str(int(self.price_start_year[3:]) + int(i / 2) + 1)),
                    trend_df
                ], how="vertical"
            )

        return trend_df.filter(
            (pl.col("DES_ZONA") == zone) &
            (pl.col("DES_STATO") == state) &
            (pl.col("DES_TIPOLOGIA") == type)
        ).select(
            "MESE_ANNO",
            "COSTO_MIN",
            "COSTO_MAX"
        )

    def get_volume_trend_mq(self, com: str, mq:int):

        # Retrieving columns name with number in the text
        allowed_mq_columns = [f for f in self.volumes_data.columns if "<" in f or ">" in f]
        mq_bounds = [int(f[1:-3]) for f in allowed_mq_columns]

        # Getting the right column base on the MQ passed in the function
        selected_mq_col = None
        for i, col in enumerate(allowed_mq_columns):
            if mq < mq_bounds[i] or i == len(allowed_mq_columns) - 1:
                selected_mq_col = col
                break

        trend_df = self._get_volume_filtered(com, str(self.volume_start_year)).select("ANNO", selected_mq_col)
        for i in range(self.volume_current_year - self.volume_start_year):
            trend_df = pl.concat([self._get_volume_filtered(com, str(self.volume_start_year + i + 1)).select("ANNO", selected_mq_col), trend_df], how="vertical")

        return trend_df

    def get_volume_market_size(self, com: str):
        return self.market_size_data.filter(pl.col("DES_COMUNE") == com).select("TAGLIA_MERCATO")

    def get_municipalities_info(self, com: str):
        return self.municipality_info.filter(pl.col("DES_COMUNE") == com)

    def get_municipalities(self):
        return self.municipality_info.select("DES_COMUNE").unique(keep='first')
