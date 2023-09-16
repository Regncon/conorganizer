import { Timestamp } from "firebase/firestore";


export type ConEvent = {
    id: string;
    title: string;
    subtitle: string;
    description: string;
    published: boolean;
    createdAt: Timestamp;
    lastUpdated: Timestamp;
    pool: string;
    gameType: string;
}