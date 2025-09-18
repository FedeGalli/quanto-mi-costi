from fastapi import FastAPI, Query
import data_manager as dm
from fastapi.middleware.cors import CORSMiddleware

data_manager = dm.DataManager()
app = FastAPI(title="Houses cost/volumes API")

# Allowed origins GO BACKEND ONLY
origins = [
    "http://localhost:5173",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,          # or ["*"] for all
    allow_credentials=True,
    allow_methods=["*"],            # or ["GET", "POST"]
    allow_headers=["*"],
)


@app.get("/get-price-volume-data")
def get_price_volume_data(
    com: str = Query(...),
    zone: str = Query(...),
    type: str = Query(...),
    state: str = Query(...),
    mq: int = Query(...)):

    current_market_size = data_manager.get_volume_market_size(com)
    current_volume_mq = data_manager.get_current_volume_mq(com, mq)
    current_min_price, current_max_price= data_manager.get_current_price(com, zone, type, state)
    volume_trend = data_manager.get_volume_trend_mq(com, mq)
    price_trend = data_manager.get_price_trend(com, zone, type, state)
    mq_range = ""
    for col in volume_trend.columns:
        if 'm²' in col:
            mq_range = col
            break


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
        "market_size": current_market_size.to_series().item(),
        "current_volume_mq": current_volume_mq.to_series().item(),
        "mq_range": mq_range,
        "current_max_price": current_max_price.to_series().item(),
        "current_min_price": current_min_price.to_series().item(),
        "volume_trend": volume_trend_CHART_DATA,
        "price_trend": price_trend_CHART_DATA
    }


@app.get("/get-municipalities-info")
def get_municipalities_info(com: str = Query(...)):
    info = data_manager.get_municipalities_info(com)
    return {
        "DATA" : info.rows(named=False)

    }

@app.get("/get-municipalities-list")
def get_municipalities_list():
    info = data_manager.get_municipalities()
    return {
        "DATA" : [t[0] for t in info.rows(named=False)]
    }
