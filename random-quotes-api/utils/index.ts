import type { QuotesDB } from "./types";

export function randomItem(arr: Array<QuotesDB>){
    if (arr.length == 0){
        return {}
    }
    return arr[Math.floor(Math.random() * arr.length)]
}