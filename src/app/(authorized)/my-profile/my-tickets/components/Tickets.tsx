import { Paper, Typography } from '@mui/material';
import Ticket from './UI/Ticket';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import BoardGameIcon from '$lib/components/icons/BoardGameIcon';
import CardGameIcon from '$lib/components/icons/CardGameIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import EnglishIcon from '$lib/components/icons/EnglishIcon';
import GamemasterIcon from '$lib/components/icons/GameMasterIcon';
import LessHoursIcon from '$lib/components/icons/LessHoursIcon';
import MiscGameIcon from '$lib/components/icons/MiscGameIcon';
import MoreHoursIcon from '$lib/components/icons/MoreHoursIcon';
import RoleplayingGameIcon from '$lib/components/icons/RoleplayingGameIcon';
import BeginnerIcon from '$lib/components/icons/BeginnerIcon';

type Props = {};

const Tickets = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = true;

    if (verifiedEmail && verifiedCheckIn) {
        return (
            <Paper sx={{ marginBottom: '2rem', paddingLeft: '2rem', width: '320px' }}>
                <Typography variant="h1">My Tickets</Typography>
                <RoleplayingGameIcon color="primary" size="large" />
                <BoardGameIcon color="primary" size="large" />
                <CardGameIcon color="primary" size="large" />
                <MiscGameIcon color="primary" size="large" />
                <GamemasterIcon color="primary" size="large" />
                <AdultsOnlyIcon color="primary" size="large" />
                <ChildFriendlyIcon color="primary" size="large" />
                <EnglishIcon color="primary" size="large" />
                <MoreHoursIcon color="primary" size="large" />
                <LessHoursIcon color="primary" size="large" />
                <BeginnerIcon color="primary" size="large" />
                <Ticket />
                <Ticket />
            </Paper>
        );
    }
    return null;
};

export default Tickets;
