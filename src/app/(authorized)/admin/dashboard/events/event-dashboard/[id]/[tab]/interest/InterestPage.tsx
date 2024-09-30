// 'use client';
import { poolTitles } from '$app/(authorized)/admin/dashboard/rooms/components/lib/helpers';
import {
    getEventById,
    getEventInterestById,
    migrateInterestsToParticipantInterests,
} from '$app/(public)/components/lib/serverAction';
import { InterestLevel, PoolName } from '$lib/enums';
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
const getOrderValue = (poolName: string): number =>
    poolName === PoolName.fridayEvening ? 1
    : poolName === PoolName.saturdayMorning ? 2
    : poolName === PoolName.saturdayEvening ? 3
    : poolName === PoolName.sundayMorning ? 4
    : 5;

const InterestPage = async ({ id }: Props) => {
    // const event = await getEventById(id);
    const eventInterests = await getEventInterestById(id);

    // const handliMoveInterests = async () => {
    //     await migrateInterestsToParticipantInterests();
    // };
    return (
        <Box>
            <Typography variant="h1">Ønskeliste:</Typography>
            {/* <Button variant="contained" color="primary">
                Migration
            </Button> */}
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit, minmax(16.7rem, 0.2fr))',
                    gap: '1rem',
                    placeContent: 'center',
                }}
            >
                {eventInterests
                    ?.sort((a, b) => {
                        const orderA = getOrderValue(a.poolName);
                        const orderB = getOrderValue(b.poolName);

                        if (orderA === orderB) {
                            // Hvis poolName er lik, sorter etter antall interesser
                            return a.interests.length - b.interests.length;
                        }
                        // Sorter etter den definerte rekkefølgen
                        return orderA - orderB;
                    })
                    ?.map((poolInterest) => (
                        <Box
                            sx={{
                                display: 'grid',
                                gap: '1rem',
                                gridTemplateRows: 'max-content repeat(3, 1fr)',
                                '*, *,  *,': {
                                    margin: 0,
                                },
                            }}
                        >
                            <Typography variant="h2">{poolTitles[poolInterest.poolName]}</Typography>
                            <PoolInterestLevel
                                interestLevel={InterestLevel.VeryInterested}
                                poolInterest={poolInterest.interests}
                            />
                            <PoolInterestLevel
                                interestLevel={InterestLevel.Interested}
                                poolInterest={poolInterest.interests}
                            />
                            <PoolInterestLevel
                                interestLevel={InterestLevel.SomewhatInterested}
                                poolInterest={poolInterest.interests}
                            />
                        </Box>
                    ))}
            </Box>
        </Box>
    );
};

export default InterestPage;
