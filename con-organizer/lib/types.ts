import { Timestamp } from 'firebase/firestore';
import { gameType, pool } from './enums';

export type ConEvent = {
    id: string;
    title: string;
    subtitle: string;
    description: string;
    published: boolean;
    createdAt: Timestamp;
    lastUpdated: Timestamp;
    pool: pool;
    gameType: gameType;
    gameSystem: string;
    room: string
};
