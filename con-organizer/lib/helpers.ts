import { FirebaseCollections } from '@/models/enums';
import { CrmJson, EventTicket, Participant, UserSettings } from '@/models/types';
import { adminDb } from './firebaseAdmin';

export function a11yProps(index: number) {
    return {
        id: `full-width-tab-${index}`,
        'aria-controls': `full-width-tabpanel-${index}`,
    };
}


