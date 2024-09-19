
import { Paper, Typography } from '@mui/material';
import Ticket from './UI/Ticket';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import RoleplayingGameIcon from '$lib/icons/roleplayinggameicon';
import BoardGameIcon from '$lib/icons/boardgameicon';
import MiscGameIcon from '$lib/icons/miscgameicon';
import CardGameIcon from '$lib/icons/cardgameicon';
import GamemasterIcon from '$lib/icons/gamemastericon';
import AdultsOnlyIcon from '$lib/icons/adultsonlyicon';
import { ChildFriendly } from '@mui/icons-material';
import ChildFriendlyIcon from '$lib/icons/childfriendlyicon';
import EnglishIcon from '$lib/icons/englishicon';
import MoreHoursIcon from '$lib/icons/morehoursicon';
import LessHoursIcon from '$lib/icons/lesshoursicon';
import BeginnerIcon from '$lib/icons/beginnericon';

type Props = {};

const Tickets = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = true;

    if (verifiedEmail && verifiedCheckIn) {
        return (
            <Paper sx={{ marginBottom: '2rem', paddingLeft: '2rem', width: '320px' }}>
                <Typography variant="h1">My Tickets</Typography>
                <RoleplayingGameIcon color='primary' size='large' />
                <BoardGameIcon color='primary' size='large' />
                <CardGameIcon color='primary' size='large' />
                <MiscGameIcon color='primary' size='large' />
                <GamemasterIcon color='primary' size='large' />
                <AdultsOnlyIcon color='primary' size='large' />
                <ChildFriendlyIcon color='primary' size='large' />
                <EnglishIcon color='primary' size='large' />
                <MoreHoursIcon color='primary' size='large' />
                <LessHoursIcon color='primary' size='large' />
                <BeginnerIcon color='primary' size='large' />
                <Ticket />
                <Ticket />
            </Paper>
        );
    }
    return null;
};

export default Tickets;
