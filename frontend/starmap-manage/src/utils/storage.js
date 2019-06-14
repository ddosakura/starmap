class StorageCache {
    constructor() {
        this.cache = {}
    }

    get(k) {
        if (this.cache[k] === undefined) {
            this.cache[k] = localStorage[k]
        }
        return this.cache[k]
    }

    set(k, v) {
        if (v !== this.cache[k]) {
            this.cache[k] = v
            localStorage[k] = v
        }
    }
}

const cache = new StorageCache()
const drivers = {}
export default drivers

class KVDriver {
    constructor(ns) {
        this.ns = ns
        this.parsers = {}
        this.cache = {}
    }

    parser(k, {
        parse = () => {},
        toString = () => {}
    }) {
        this.parsers[this.ns + "-" + k] = {
            parse,
            toString,
        }
    }

    _getParser(k) {
        const parser = this.parsers[k]
        if (parser !== undefined && parser.parse instanceof Function && parser.toString instanceof Function) {
            return parser
        }
        return {}
    }

    get(k) {
        k = this.ns + "-" + k
        if (this.cache[k] === undefined) {
            const {
                parse
            } = this._getParser(k)
            if (!parse) {
                return [null, ErrUnknowParser]
            }
            const d = cache.get(k)
            if (d === undefined) {
                return [null, ErrNotExist]
            }
            this.cache[k] = parse(d)
        }
        return [this.cache[k], null]
    }

    set(k, v) {
        k = this.ns + "-" + k
        const {
            toString
        } = this._getParser(k)
        if (!toString) {
            return [null, ErrUnknowParser]
        }
        this.cache[k] = v
        if (v !== undefined) {
            v = toString(v)
        }
        cache.set(k, v)
        return null
    }
}

export const JSONParse = {
    parse: JSON.parse,
    toString: JSON.stringify,
}
export const ErrUnknowParser = "ErrUnknowParser"
export const ErrNotExist = "ErrNotExist"

drivers.cache = cache
drivers.kv = new KVDriver("kv")