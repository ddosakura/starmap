import drivers, {
    JSONParse,
    ErrNotExist,
} from './storage'

const {
    kv
} = drivers
kv.parser("jwt", JSONParse)

export function getUserInfo() {
    const [d, e] = kv.get("jwt")
    if (e === ErrNotExist) {
        return [ undefined, false ]
    }
    const t = d.payload.exp - new Date() / 1000
    return [d.payload.UserInfo, t > 0]
}

export function getToken() {
    const [d, e] = kv.get("jwt")
    if (e === ErrNotExist) {
        return
    }
    if (d.payload.exp - new Date() / 1000 <= 0) {
        return
    }
    return d.token
}

export function freshJWT(token) {
    const b = new Buffer(token.split(".")[1], 'base64')
    const s = b.toString()
    const payload = JSON.parse(s)

    kv.set("jwt", {
        token,
        payload,
    })
}

export function cleanJWT() {
    kv.set("jwt", undefined)
}
