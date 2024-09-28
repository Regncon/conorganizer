'use client';

import { Button, Box, Slider, sliderClasses, Typography } from '@mui/material';
import Link from 'next/link';
import HelpIcon from '@mui/icons-material/Help';
import { use, useCallback, useEffect, useState, useTransition } from 'react';
import Image from 'next/image';
import AwakeDragons from 'public/interessedragene/2024AwakeDragons1_1.png';
import HappyDragons from 'public/interessedragene/2024HappyDragons1_1.png';
import SleepyDragons from 'public/interessedragene/2024SleepyDragons1_1.png';
import VeryHappyDragons from 'public/interessedragene/2024VeryHappyDragons1_1.png';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { InterestLevel, PoolName } from '$lib/enums';
import { ParticipantCookie } from '$lib/types';
import { getInterest, updateInterest } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import debounce from '$lib/debounce';

const poolTitlesWithTime = {
    [PoolName.fridayEvening]: 'Fredag Kveld Kl 18 - 23',
    [PoolName.saturdayMorning]: 'Lørdag Morgen Kl 10 - 15',
    [PoolName.saturdayEvening]: 'Lørdag Kveld Kl 18 - 23',
    [PoolName.sundayMorning]: 'Søndag Morgen Kl 10 - 15',
};
type Props = {
    poolName?: PoolName;
    poolEventId?: string;
    poolEventTitle?: string;
    disabled?: boolean;
    activeParticipant?: { id?: string; interestLevel?: InterestLevel };
};

const InterestSelector = ({ poolName, poolEventId, poolEventTitle, disabled = false, activeParticipant }: Props) => {
    const [interest, setInterest] = useState<number>(
        typeof activeParticipant?.interestLevel === 'string' ? 0 : (activeParticipant?.interestLevel ?? 0)
    );
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const notLoggedIn = activeParticipant?.id === undefined && activeParticipant?.interestLevel === undefined;

    const handleParticipantSelect = (e: Event & { detail?: { loading?: boolean } }) => {
        const isLoading = e.detail?.loading;
        setIsLoading(isLoading ?? false);
    };

    useEffect(() => {
        console.log('initial interestLevel: ', activeParticipant?.interestLevel);
        setInterest(typeof activeParticipant?.interestLevel === 'string' ? 0 : (activeParticipant?.interestLevel ?? 0));
    }, [activeParticipant]);
    useEffect(() => {
        document.addEventListener('my-participants-changed', handleParticipantSelect);
        return () => {
            document.removeEventListener('my-participants-changed', handleParticipantSelect);
        };
    }, []);

    const [isPending, startTransition] = useTransition();
    const startTransitionDebounced = useCallback(debounce(startTransition, 500), []);

    const incrementInterestButton = () => {
        let interestLevel: InterestLevel = InterestLevel.NotInterested;
        if (interest === 3) {
            interestLevel = InterestLevel.NotInterested;
            setInterest(0);
        } else {
            interestLevel = interest + 1;
            setInterest(interestLevel);
        }

        if (activeParticipant?.id && poolEventId) {
            startTransitionDebounced(async () => {
                console.log('incrementInterestButton: ', interestLevel);

                await updateInterest(activeParticipant?.id, poolEventId, poolName, poolEventTitle, interestLevel);
            });
        }
    };
    const incrementInterestSlider = (interest: InterestLevel) => {
        if (activeParticipant?.id && poolEventId) {
            startTransitionDebounced(async () => {
                console.log('incrementInterestSlider: ', interest);

                await updateInterest(activeParticipant?.id, poolEventId, poolName, poolEventTitle, interest);
            });
        }
    };

    const marks = [
        {
            value: 1,
            label: (
                <Box sx={{ display: 'flex', width: '25rem', alignItems: 'center', gap: '0.5rem' }}>
                    <Image src={SleepyDragons} alt="Ikke interessert" width={100} height={60} />
                    <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                        {poolName ?
                            <Typography>{poolTitlesWithTime[poolName]}</Typography>
                            : null}
                        <Typography>Ikke interessert</Typography>
                    </Box>
                </Box>
            ),
        },
        {
            value: 2,
            label: (
                <Box sx={{ display: 'flex', width: '25rem', alignItems: 'center', gap: '0.5rem' }}>
                    <Image src={AwakeDragons} alt="Litt interessert" width={100} height={60} />
                    <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                        {poolName ?
                            <Typography>{poolTitlesWithTime[poolName]}</Typography>
                            : null}
                        <Typography>Litt interessert</Typography>
                    </Box>
                </Box>
            ),
        },
        {
            value: 3,
            label: (
                <Box sx={{ display: 'flex', width: '25rem', alignItems: 'center', gap: '0.5rem' }}>
                    <Image src={HappyDragons} alt="Interessert" width={100} height={60} />
                    <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                        {poolName ?
                            <Typography>{poolTitlesWithTime[poolName]}</Typography>
                            : null}
                        <Typography>Interessert</Typography>
                    </Box>
                </Box>
            ),
        },
        {
            value: 4,
            label: (
                <Box sx={{ display: 'flex', width: '25rem', alignItems: 'center', gap: '0.5rem' }}>
                    <Image src={VeryHappyDragons} alt="Veldig interessert" width={100} height={60} />
                    <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                        {poolName ?
                            <Typography>{poolTitlesWithTime[poolName]}</Typography>
                            : null}
                        <Typography>Veldig interessert</Typography>
                    </Box>
                </Box>
            ),
        },
    ];

    return (
        <>
            <ParticipantSelector poolEventId={poolEventId} />
            <Button
                variant="contained"
                color="primary"
                fullWidth
                sx={{
                    fontSize: '1.2rem',
                    textTransform: 'none',
                    marginBlock: '1rem',
                    minHeight: '62px',
                    maxWidth: 'var(--slider-interest-width)',
                }}
                onClick={incrementInterestButton}
                disabled={disabled || isPending || notLoggedIn || isLoading}
            >
                {isPending ? 'lagrer :)' : marks[interest].label}
            </Button>
            <Box sx={{ padding: '0.35rem', marginBottom: '0.1rem', maxWidth: 'var(--slider-interest-width)' }}>
                <Slider
                    onChange={(e) => {
                        const target = e.target as HTMLInputElement;
                        setInterest(Number(target.value));
                        incrementInterestSlider(Number(target.value));
                    }}
                    disabled={disabled || isPending || notLoggedIn || isLoading}
                    sx={{
                        color: 'primary.main',
                        [`.${sliderClasses.rail}`]: { backgroundColor: '#3d3b3b', height: '1rem' },
                        [`.${sliderClasses.track}`]: { height: '1rem' },
                        [`.${sliderClasses.mark}`]: {
                            borderRadius: '50%',
                            outlineColor: 'primary.main',
                            outlineWidth: '0.8rem',
                            outlineStyle: 'solid',
                            outlineOffset: '-1px',
                            opacity: '1',
                        },
                        [`.${sliderClasses.markActive}`]: { backgroundColor: 'primary.main' },
                        [`.${sliderClasses.thumb}:before`]: { boxShadow: 'unset' },
                        [`.${sliderClasses.valueLabelCircle}`]: { display: 'none' },
                    }}
                    marks
                    value={interest}
                    min={0}
                    max={3}
                />
            </Box>
            <Typography sx={{ marginBottom: '0.8rem', fontSize: '1rem' }}>
                Dra baren over for å melde din interesse!
            </Typography>
            <Box>
                <Box sx={{ marginBottom: '3rem' }}>
                    <Box
                        component={Link}
                        href="/hjelppaamelding"
                        sx={{
                            display: 'inline-flex',
                            gap: '0.4rem',
                            paddingLeft: '0.5rem',
                            color: 'primary.main',
                            maxWidth: 'var(--slider-interest-width)',
                        }}
                    >
                        <HelpIcon sx={{ scale: '1.5', placeSelf: 'center' }} />
                        <Typography component="p">Forvirret? Les mer om påmeldingsystemet</Typography>
                    </Box>
                </Box>
            </Box>
        </>
    );
};

export default InterestSelector;
