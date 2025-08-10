from fastapi import FastAPI
import data_manager as dm
from pydantic import BaseModel
from fastapi import Query

data_manager = dm.DataManager()
app = FastAPI(title="Houses cost/volumes API")


# Example request schema
class DataRequest(BaseModel):
    com: str
    zone: str
    type: str
    state: str
    mq: int

@app.get("/get-price-volume-data")
def get_price_volume_data(com: str = Query(...),
    zone: str = Query(...),
    type: str = Query(...),
    state: str = Query(...),
    mq: int = Query(...)):

    current_volume_mq = data_manager.get_current_volume_mq(com, mq)
    current_max_price, current_min_price= data_manager.get_current_price(com, zone, type, state)
    volume_trend = data_manager.get_volume_trend_mq(com, mq)
    price_trend = data_manager.get_price_trend(com, zone, type, state)

    print(current_volume_mq)

    mq_volume_column = volume_trend.columns[1]
    volume_trend_CHART_DATA = {
        "labels": volume_trend["ANNO"].to_list(),
        "datasets": [
            {
                "label": "Volume",
                "data": volume_trend[mq_volume_column].to_list()
            }
        ]
    }

    price_trend_CHART_DATA = {
        "labels": price_trend["MESE_ANNO"].to_list(),
        "datasets": [
            {
                "label": "Costo min",
                "data": price_trend["COSTO_MIN"].to_list()
            },
            {
                "label": "Costo max",
                "data": price_trend["COSTO_MAX"].to_list()
            }
        ]
    }
    # Return processed data as list of dicts
    return {
        "current_volume_mq": current_volume_mq.to_series().item(),
        "current_max_price": current_max_price.to_series().item(),
        "current_min_price": current_min_price.to_series().item(),
        "volume_trend": volume_trend_CHART_DATA,
        "price_trend": price_trend_CHART_DATA
    }
