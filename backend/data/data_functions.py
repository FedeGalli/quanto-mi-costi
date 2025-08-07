from posix import stat
from turtledemo.forest import start
import polars as pl
from data_aggregator import get_volumes_market_size_df, get_volumes_df, get_prices_df

class DataManager:

    def __init__(self):
        self.start_year = 2018
        self.volumes_data = get_volumes_df()
        self.prices_data = get_prices_df()
        self.market_size_data = get_volumes_market_size_df()

    def get_volumes_market_size(self, com: str):
        return self.market_size_data.filter(pl.col("DES_COMUNE") == com).select("TAGLIA_MERCATO")

    def _get_volumes_filtered(self, com: str, year:str):
        return self.volumes_data.filter((pl.col("DES_COMUNE") == com) & (pl.col("ANNO") == year))

    def _get_prices_filtered(self, com: str, month_year:str):
        return self.prices_data.filter((pl.col("DES_COMUNE") == com) & (pl.col("MESE_ANNO") == month_year))

    def get_volume_mq(self, com: str, mq:int, year:str):

        # Retrieving columns name with number in the text
        allowed_mq_columns = [f for f in self.volumes_data.columns if "<" in f or ">" in f]
        mq_bounds = [int(f[1:-3]) for f in allowed_mq_columns]

        selected_mq_col = None
        for i, col in enumerate(allowed_mq_columns):
            if mq < mq_bounds[i] or i == len(allowed_mq_columns) - 1:
                selected_mq_col = col
                break

        return self._get_volumes_filtered(com, year).select(selected_mq_col)

    def get_min_price(self, com:str, month_year:str, zone:str, type: str, state:str):
        return self._get_prices_filtered(com, month_year).filter(
            (pl.col("DES_ZONA") == zone) &
            (pl.col("DES_STATO") == state) &
            (pl.col("DES_TIPOLOGIA") == type)
        ).select("COSTO_MIN")

    def get_max_price(self, com:str, month_year:str, zone:str, type: str, state:str):
        return self._get_prices_filtered(com, month_year).filter(
            (pl.col("DES_ZONA") == zone) &
            (pl.col("DES_STATO") == state) &
            (pl.col("DES_TIPOLOGIA") == type)
        ).select("COSTO_MAX")
    def get_price_trend(self):
        pass
    def get_volume_trend(self, com: str, year:str):
        trend_df = self._get_volumes_filtered(com, str(self.start_year))
        for i in range(2024 - self.start_year):
            trend_df = pl.concat([self._get_volumes_filtered(com, str(self.start_year + i + 1)), trend_df], how="vertical")

        return trend_df

d = DataManager()
print(d.get_volume_trend("RONCADELLE", "2024"))
