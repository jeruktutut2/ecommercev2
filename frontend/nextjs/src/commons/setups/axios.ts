import axios, { AxiosInstance } from "axios";

let httpClient: AxiosInstance | null = null

export function getHttpClient() {
    if (!httpClient) {
        httpClient = axios.create({
            baseURL: "/",
            timeout: 60000
        })
    }
    return httpClient
}