import { Timestamp } from 'firebase/firestore';
import { EnrollmentChoice, GameType, Pool } from '@/lib/enums';

export type ConEvent = {
    id: string;
    title: string;
    subtitle: string;
    description: string;
    published: boolean;
    imageUrl: string;
    createdAt: Timestamp;
    lastUpdated: Timestamp;
    pool: Pool;
    gameType: GameType;
    gameSystem: string;
    room: string;
    host: string;
};


export type UserSettings = {
    id: string;
    admin: boolean;
    name: string;
};

export type UserData = {
    id: string;
    enrollments: Enrollment[];
    createdAt: Timestamp;
    lastUpdated: Timestamp;
};

export type Enrollment = {
    id: string;
    userId: string;
    eventId: string;
    choice: EnrollmentChoice;  
    createdAt: Timestamp;
    lastUpdated: Timestamp;
};