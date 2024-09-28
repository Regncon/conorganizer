// 'use client';
import { poolTitles } from '$app/(authorized)/admin/dashboard/rooms/components/lib/helpers';
import {
    getEventById,
    getEventInterestById,
    migrateInterestsToParticipantInterests,
} from '$app/(public)/components/lib/serverAction';
import { InterestLevel } from '$lib/enums';
import { Interest } from '$lib/types';
import { Box, Button, Link, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import NextLink from 'next/link';
import AwakeDragons from 'public/interessedragene/2024AwakeDragons1_1.png';
import HappyDragons from 'public/interessedragene/2024HappyDragons1_1.png';
import SleepyDragons from 'public/interessedragene/2024SleepyDragons1_1.png';
import VeryHappyDragons from 'public/interessedragene/2024VeryHappyDragons1_1.png';
import PoolInterestLevel from './components/PoolInterestLeve';

type Props = {
    id: string;
};

const InterestPage = async ({ id }: Props) => {
    // const event = await getEventById(id);
    const eventInterests = await getEventInterestById(id);

    // const handliMoveInterests = async () => {
    //     await migrateInterestsToParticipantInterests();
    // };
    return (
        <Box>
            <Typography variant="h1">Ønskeliste:</Typography>
            <Button variant="contained" color="primary">
                Migration
            </Button>
            {eventInterests.map((poolInterest) => (
                <Paper>
                    <Typography variant="h2">{poolTitles[poolInterest.poolName]}</Typography>
                    <PoolInterestLevel
                        interestLevel={InterestLevel.VeryInterested}
                        poolInterest={poolInterest.interests}
                    />
                    <PoolInterestLevel interestLevel={InterestLevel.Interested} poolInterest={poolInterest.interests} />
                    <PoolInterestLevel
                        interestLevel={InterestLevel.SomewhatInterested}
                        poolInterest={poolInterest.interests}
                    />
                </Paper>
            ))}
        </Box>
    );
};

export default InterestPage;
