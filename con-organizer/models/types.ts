import { Timestamp } from 'firebase/firestore';
import { GameType, Pool } from '@/lib/enums';

export type ConEvent = {
    id: string;
    title: string;
    subtitle: string;
    description: string;
    published: boolean;
    createdAt: Timestamp;
    lastUpdated: Timestamp;
    pool: Pool;
    gameType: GameType;
    gameSystem: string;
    room: string;
    host: string;
};
