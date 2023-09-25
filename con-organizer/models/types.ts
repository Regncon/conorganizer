import { Timestamp } from 'firebase/firestore';
import { EnrollmentChoice, GameType, Pool } from '@/models/enums';

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
    hideEnrollment: boolean;
    sortingIndex: number;
    childFriendly: boolean;
    possiblyEnglish: boolean;
    adultsOnly: boolean;
    volunteersPossible: boolean;
    lessThanThreeHours: boolean;
    moreThanSixHours: boolean;
    beginnerFriendly: boolean;
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
    choice: EnrollmentChoice;
    createdAt: Timestamp;
    lastUpdated: Timestamp;
};
