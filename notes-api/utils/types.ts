export interface Note {
    id: string;
    title: string;
    body: string;
    created_at: number;
    updated_at: number;
}

export interface UpdateNoteRequest{
    id: string;
    title?: string;
    body?: string;
}