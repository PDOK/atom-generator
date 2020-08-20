def build_uri(*parts, ending_slash=False):
    """
    Combines elements to a uri with slashes and an optional ending slash.

    Args:
        *parts: variable length list of strings
        ending_slash (bool): adds a `/` to the end when True. Defaults to False.

    Returns:
        str: formatted uri
    """
    return "/".join(p.strip("/") for p in parts if p) + "/" * ending_slash
