import { type } from 'os';
import { PoolName, type RoomName } from './enums';

export type EventCardProps = {
    title: string;
    gameMaster: string;
    system: string;
    shortDescription: string;
    icons?: string[];
    backgroundImage?: string;
    myEventBar?: boolean;
    myEventBarSubmitted?: boolean;
    myEventDocId?: string;
    isAccepted?: boolean;
};

export type ConEvent = {
    published: boolean;
    id?: string;
    title: string;
    gameMaster: string;
    system: string;
    shortDescription: string;
    description: string;
    smallImageURL?: string;
    bigImageURL?: string;
    icons?: string[];
    email: string;
    name: string;
    phone: string;
    gameType: string;
    isSmallCard: boolean;
    participants: number;
    puljeFridayEvening: boolean;
    puljeSaturdayMorning: boolean;
    puljeSaturdayEvening: boolean;
    puljeSundayMorning: boolean;
    poolIds: PoolChildRef[];
    roomIds: RoomChildRef[];
    unwantedFridayEvening: boolean;
    unwantedSaturdayMorning: boolean;
    unwantedSaturdayEvening: boolean;
    unwantedSundayMorning: boolean;
    moduleCompetition: boolean;
    childFriendly: boolean;
    possiblyEnglish: boolean;
    adultsOnly: boolean;
    volunteersPossible: boolean;
    lessThanThreeHours: boolean;
    moreThanSixHours: boolean;
    beginnerFriendly: boolean;
    additionalComments: string;
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
    subTitle: string;
};

export type RoomPlayer = {
    email: string;
    name: string;
    id: string;
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
};
export type EventRoom = {
    id?: string;
    name: RoomName;
    eventId: string;
    poolId: string;
    players: RoomPlayer[];
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
};

export type RoomChildRef = {
    id: string;
    poolId: string;
    poolName: PoolName;
    roomName: RoomName;
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
};

export type RoomItemInfo = {
    roomName: RoomName;
    top: number;
    left: number;
};

export type PoolChildRef = {
    id: string;
    poolName: PoolName;
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
};
export type IconTypes = {
    childFriendly: boolean;
    possiblyEnglish: boolean;
    adultsOnly: boolean;
    lessThanThreeHours: boolean;
    moreThanSixHours: boolean;
    beginnerFriendly: boolean;
};

export type IconName = keyof IconTypes;
export type IconOption = {
    label: string;
    iconName: IconName;
};
export type PoolEvent = {
    poolName: PoolName;
    published: boolean;
    id?: string;
    parentEventId: string;
    title: string;
    gameMaster: string;
    system: string;
    shortDescription: string;
    description: string;
    smallImageURL?: string;
    bigImageURL?: string;
    gameType: string;
    isSmallCard: boolean;
    participants: number;
    childFriendly: boolean;
    possiblyEnglish: boolean;
    adultsOnly: boolean;
    volunteersPossible: boolean;
    lessThanThreeHours: boolean;
    moreThanSixHours: boolean;
    beginnerFriendly: boolean;
    additionalComments: string;
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
    icons?: IconOption[];
};

export type MyNewEvent = {
    id: string;
    eventDocId?: string;
    email: string;
    name: string;
    phone: string;
    title: string;
    system: string;
    gameType: string;
    participants: number;
    fridayEvening: boolean;
    saturdayMorning: boolean;
    saturdayEvening: boolean;
    sundayMorning: boolean;
    moduleCompetition: boolean;
    childFriendly: boolean;
    possiblyEnglish: boolean;
    adultsOnly: boolean;
    volunteersPossible: boolean;
    lessThanThreeHours: boolean;
    moreThanSixHours: boolean;
    beginnerFriendly: boolean;
    description: string;
    additionalComments: string;
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
    subTitle: string;
    isSubmitted: boolean;
    isRead?: boolean;
    isAccepted?: boolean;
};
export type Participant = {
    id: string;
    name: string;
    users?: string[];
    over18: boolean;
    orderId: number;
    ticketId: number;
    ticketEmail: string;
    oredrEmails: string[];
    ticketCategory: string;
    ticketCategoryId: number;
    connectedEmails: string[];
    createdAt: string;
    createdBy: string;
    updateAt: string;
    updatedBy: string;
};

export type ConUser = {
    id?: string;
    email?: string;
    admin?: boolean;
    participants?: Participant[];
};

export type ActionResponse = {
    type: 'info' | 'success' | 'warning' | 'error';
    error?: string;
    message: string;
};
