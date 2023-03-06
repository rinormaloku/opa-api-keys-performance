package apikeys

default allow = false

allow {
    has_key(data.ApiKeys, input.apikey)
}

deny {
	not has_key(data.ApiKeys, input.apikey)
}

key_data = data.ApiKeys[input.apikey]

has_key(x, k) {
	x[k]
}