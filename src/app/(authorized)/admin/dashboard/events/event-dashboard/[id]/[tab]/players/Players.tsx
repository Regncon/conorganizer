import {
    Box,
    Card,
    CardActions,
    CardHeader,
    FormControlLabel,
    FormGroup,
    Link,
    Paper,
    Switch,
    Typography,
} from '@mui/material';
import PlayerInterestInfo from './components/PlayerInterestInfo';
import { PlayerInterest } from '$lib/types';
import { InterestLevel, PoolName } from '$lib/enums';

type Props = {
    id: string;
};

const Players = async ({ id }: Props) => {
    console.log('Players', id);

    const dummyPlayerInrerestData: PlayerInterest = {
        interestLevel: InterestLevel.VeryInterested,
        poolEventId: '',
        participantId: '',
        firstName: 'Kari',
        lastName: 'Nordmann',
        isOver18: true,
        ticketCategoryID: 0,
        ticketCategory: 'Festivalpass Ungdom/student (13-25/30år) Early-bird',
        conPlayers: [
            {
                id: undefined,
                participantId: '',
                firstName: '',
                lastName: '',
                interestLevel: InterestLevel.VeryInterested,
                poolEventId: '',
                poolEventTitle: 'Et gøyalt spill',
                poolName: PoolName.fridayEvening,
                isFirstChoice: true,
                isGameMaster: false,
                createdAt: '',
                createdBy: '',
                updateAt: '',
                updatedBy: '',
            },
            {
                id: undefined,
                participantId: '',
                firstName: '',
                lastName: '',
                interestLevel: InterestLevel.SomewhatInterested,
                poolEventId: '',
                poolEventTitle: '"Hva f**n skjedde!?!??" -- ghost/echo av John Harper',
                poolName: PoolName.saturdayMorning,
                isFirstChoice: false,
                isGameMaster: true,
                createdAt: '',
                createdBy: '',
                updateAt: '',
                updatedBy: '',
            },
            {
                id: undefined,
                participantId: '',
                firstName: '',
                lastName: '',
                interestLevel: InterestLevel.Interested,
                poolEventId: '',
                poolEventTitle: 'Random dnd event',
                poolName: PoolName.saturdayMorning,
                isFirstChoice: false,
                isGameMaster: true,
                createdAt: '',
                createdBy: '',
                updateAt: '',
                updatedBy: '',
            },
        ],
    };
    return (
        <Paper sx={{ paddingBlock: '1rem' }}>
            <Typography variant="h1">Spillere:</Typography>
            <Card sx={{ backgroundColor: 'rgb(55, 59, 87)' }}>
                <CardHeader title="Prioriterg" />
                <CardActions>
                    <FormGroup>
                        <FormControlLabel control={<Switch />} label="Under 18" labelPlacement="start" />
                        <FormControlLabel control={<Switch />} label="Over 18" labelPlacement="start" />
                        <FormControlLabel
                            control={<Switch />}
                            label="SØNDAG Dagspass Barn (3-12)"
                            labelPlacement="start"
                        />
                    </FormGroup>
                </CardActions>
            </Card>
            <PlayerInterestInfo playerInterest={dummyPlayerInrerestData} />
        </Paper>
    );
};

export default Players;
