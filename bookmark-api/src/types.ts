export interface UrlDB {
    id: string
    url: string;
    title: string;
    tags?: Array<string>;
}

export interface getreq{
    tag?: string
}