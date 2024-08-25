import { geMyEventByRefPath } from '$app/(authorized)/my-events/actions';
import MainEvent from '$app/(public)/event/[id]/event';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Event, MyNewEvent } from '$lib/types';
import { Box, Typography } from '@mui/material';

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

    const event: Event = {
        id: nyEvent.id,
        name: nyEvent.name,
        title: nyEvent.title,
        subTitle: nyEvent.subTitle,
        description: nyEvent.description,
        gameType: nyEvent.gameType,
        system: nyEvent.system,
        data: false,
        gameMaster: nyEvent.name,
        shortDescription: '',
        email: '',
        phone: '',
        participants: 0,
        fridayEvening: false,
        saturdayMorning: false,
        saturdayEvening: false,
        sundayMorning: false,
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
