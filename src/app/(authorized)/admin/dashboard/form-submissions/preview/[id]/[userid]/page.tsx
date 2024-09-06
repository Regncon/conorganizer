import { geMyEventByRefPath } from '$app/(authorized)/my-events/actions';
import MainEvent from '$app/(public)/event/[id]/MainEvent';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ConEvent, MyNewEvent } from '$lib/types';
import { Box, Typography } from '@mui/material';
import { updateReadAndOrAcceptedStatus } from '../../../actions';

type Props = {
    params: {
        id: string;
        userid: string;
    };
};
const FormSubmissionsPreviewPage = async ({ params: { id, userid } }: Props) => {
    const { db } = await getAuthorizedAuth();
    if (db === null) {
        return;
    }
    const nyEvent: MyNewEvent = await geMyEventByRefPath(db, id, userid);

    if (nyEvent.isRead === false) {
        updateReadAndOrAcceptedStatus(`users/${userid}/my-events/${id}`, {
            isRead: true,
        });
    }

    const event: ConEvent = {
        id: nyEvent.id,
        name: nyEvent.name,
        title: nyEvent.title,
        subTitle: nyEvent.subTitle,
        description: nyEvent.description,
        gameType: nyEvent.gameType,
        system: nyEvent.system,
        gameMaster: nyEvent.name,
        published: false,
        shortDescription: '',
        email: '',
        phone: '',
        isSmallCard: false,
        participants: 0,
        puljeFridayEvening: false,
        puljeSaturdayMorning: false,
        puljeSaturdayEvening: false,
        puljeSundayMorning: false,
        unwantedFridayEvening: false,
        unwantedSaturdayMorning: false,
        unwantedSaturdayEvening: false,
        unwantedSundayMorning: false,
        moduleCompetition: false,
        childFriendly: false,
        possiblyEnglish: false,
        adultsOnly: false,
        volunteersPossible: false,
        lessThanThreeHours: false,
        moreThanSixHours: false,
        beginnerFriendly: false,
        additionalComments: '',
        createdAt: '',
        createdBy: '',
        updateAt: '',
        updatedBy: '',
    };
    return (
        <Box sx={{ maxWidth: '430px', margin: 'auto' }}>
            <Typography variant="h1">Forhåndsvisning</Typography>
            <hr />
            <MainEvent eventData={event} />
        </Box>
    );
};
export default FormSubmissionsPreviewPage;
