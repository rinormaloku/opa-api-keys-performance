package apikeys

allow {
    has_key(data.ApiKeys, input.apikey)
}

deny {
	not has_key(data.ApiKeys, input.apikey)
}

has_key(x, k) {
	x[k]
}