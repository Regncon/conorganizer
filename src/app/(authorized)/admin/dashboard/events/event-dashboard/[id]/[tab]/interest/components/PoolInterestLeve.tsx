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

type Props = {
    interestLevel: InterestLevel;
    poolInterest: Interest[];
};

const InterestLevelToLable = {
    [InterestLevel.VeryInterested]: 'Veldig interessert',
    [InterestLevel.Interested]: 'Interessert',
    [InterestLevel.SomewhatInterested]: 'Litt interessert',
    [InterestLevel.NotInterested]: 'Ikke interessert',
};

const interestLevelToImage = {
    [InterestLevel.VeryInterested]: VeryHappyDragons,
    [InterestLevel.Interested]: HappyDragons,
    [InterestLevel.SomewhatInterested]: SleepyDragons,
    [InterestLevel.NotInterested]: AwakeDragons,
};

const PoolInterestLevel = ({ interestLevel, poolInterest }: Props) => {
    return (
        <Box>
            <Typography variant="h3">{InterestLevelToLable[interestLevel]}</Typography>

            <Image
                src={interestLevelToImage[interestLevel]}
                alt={InterestLevelToLable[interestLevel]}
                width={100}
                height={50}
            />
            {poolInterest
                .filter((interest) => interest.interestLevel === interestLevel)
                .map((veryInterested) => (
                    <Box>
                        <Typography variant="h3">{veryInterested.participantFirstName}</Typography>
                    </Box>
                ))}
        </Box>
    );
};
export default PoolInterestLevel;
