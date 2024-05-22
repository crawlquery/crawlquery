
export interface Page {
    id: string;
    url: string;
    title: string;
    meta_description: string;
    content: string;
}

export interface Result {
    page_id: string;
    score: number;
    page: Page;
}

export interface Node {
    id: string;
    key: string;
    hostname: string;
    port: number;
    shard_id: number;
    created_at: string;
}