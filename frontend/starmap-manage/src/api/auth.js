import {
    _axios
} from './index'

export async function getInfo() {
    const res = await _axios.post("/auth/info")
    return res.data
}

export async function login(user, pass) {
    const res = await _axios.post("/auth/login", {
        user,
        pass,
    })
    return res.data
}