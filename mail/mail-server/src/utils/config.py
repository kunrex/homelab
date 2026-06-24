import yaml

from ..datatypes import Config

def load_config(path: str = "config.yaml") -> Config | None:
    try:
        with open(path) as f:
            return Config(**yaml.safe_load(f))
    except Exception:
        return None
