'use client';
import { Box, CircularProgress, Paper, TextField, Typography, useTheme, type SxProps, type Theme } from '@mui/material';
import Image from 'next/image';
import diceSmall from '$public/dice-small.webp';
import { useEffect, useState } from 'react';
import { faUserSecret } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faScroll } from '@fortawesome/free-solid-svg-icons/faScroll';
import { type ConEvent, type PoolEvent } from '$lib/types';
import { db } from '$lib/firebase/firebase';
import { onSnapshot, doc, type Unsubscribe } from 'firebase/firestore';
import { MuiMarkdown } from 'mui-markdown';
import NavigatePreviousLink from './ui/NavigatePreviousLink';
import NavigateNextLink from './ui/NavigateNextLink';
import BackButton from '$app/(authorized)/components/BackButton';
import InterestSelector from './components/InterestSelector';
import { createIconArray } from '$app/(public)/components/lib/helpers/icons';
import ChipCarousel from './ui/ChipCarousel';
import GoToEventAdministrationButton from './ui/GoToEventAdministrationButton';

export const dynamic = 'force-static';

type Props = {
    id?: string;
    eventData?: PoolEvent;
    editable?: boolean;
    editDescription?: (edit: boolean) => void;
    prevNavigationId?: string;
    nextNavigationId?: string;
    parent?: boolean;
    handleChange?: (data: Partial<ConEvent>) => Promise<void>;
    isAdmin?: boolean;
};

const MainEvent = ({
    id,
    eventData,
    editable = false,
    parent = false,
    editDescription,
    prevNavigationId,
    nextNavigationId,
    handleChange,
    isAdmin = false,
}: Props) => {
    const [data, setData] = useState<PoolEvent | undefined>(eventData);
    const [isEditingTitle, setIsEditingTitle] = useState<boolean>(false);
    const [isEditingGameMaster, setIsEditingGameMaster] = useState<boolean>(false);
    const [isEditingSystem, setIsEditingSystem] = useState<boolean>(false);
    const [isEditingShortDescription, setIsEditingShortDescription] = useState<boolean>(false);

    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (id !== undefined && eventData === undefined) {
            unsubscribeSnapshot = onSnapshot(doc(db, `${parent ? 'events' : 'pool-events'}`, id), (snapshot) => {
                let poolEvent = snapshot.data() as PoolEvent | undefined;
                let icons = createIconArray({
                    adultsOnly: poolEvent?.adultsOnly ?? false,
                    childFriendly: poolEvent?.childFriendly ?? false,
                    beginnerFriendly: poolEvent?.beginnerFriendly ?? false,
                    lessThanThreeHours: poolEvent?.lessThanThreeHours ?? false,
                    moreThanSixHours: poolEvent?.moreThanSixHours ?? false,
                    possiblyEnglish: poolEvent?.possiblyEnglish ?? false,
                });

                if (poolEvent) {
                    poolEvent.icons = icons;
                    setData(poolEvent);
                }
            });
        }
        return () => {
            unsubscribeSnapshot?.();
        };
    }, [id, eventData]);

    const {
        palette: {
            background: { default: themeBackgroundColor },
        },
    } = useTheme();

    const paragraphStyle: SxProps<Theme> = { margin: '1rem 0' };
    return data ?
        <Paper
            elevation={1}
            sx={{
                maxWidth: 'var(--max-width)',
                margin: 'auto',
                '&': {
                    '--max-width': '430px',
                    '--arrow-size': '2rem',
                },
            }}
        >
            <Box
                sx={{
                    display: 'grid',
                    '& > *': { gridColumn: '1 / 2', gridRow: '1 / 2' },
                    '& > img': { width: '100%', height: '100%' },
                }}
            >
                <img
                    src={data?.smallImageURL ?? diceSmall}
                    width={320}
                    height={273}
                    sizes="100vw"
                    alt="noe alt-tekst"
                    loading="lazy"
                />
                <Box
                    sx={{
                        background: `linear-gradient(0deg, ${themeBackgroundColor}, transparent)`,
                    }}
                >
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateRows: '2rem 1fr',
                            height: '100%',
                            wordBreak: 'break-word',
                        }}
                    >
                        <Box sx={{ placeSelf: 'start' }}>
                            <BackButton />
                        </Box>
                        {isEditingTitle ?
                            <TextField
                                name="title"
                                value={data.title}
                                onChange={(e) => setData({ ...data, title: e.target.value })}
                                onBlur={() => setIsEditingTitle(false)}
                                autoFocus
                                variant="outlined"
                                fullWidth
                                sx={{ placeSelf: 'end center', paddingBottom: '2.5rem' }}
                            />
                            : <Typography
                                variant="h1"
                                align="center"
                                sx={{ corsor: 'pointer', placeSelf: 'end center', paddingBottom: '2.5rem' }}
                                onClick={() => editable && setIsEditingTitle(true)}
                            >
                                {data.title || 'Tittel'}
                            </Typography>
                        }
                    </Box>
                </Box>
            </Box>
            <Box sx={{ padding: '1rem' }}>
                <Box sx={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
                    <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                        <Box
                            component={FontAwesomeIcon}
                            icon={faUserSecret}
                            size="2x"
                            sx={{ color: 'primary.main' }}
                        />
                        <Box>
                            <Typography component="span" sx={{ color: 'primary.main' }}>
                                {data.gameType === 'rolePlaying' ? 'Gamemaster' : 'Arrangør'}
                            </Typography>
                            {isEditingGameMaster ?
                                <TextField
                                    name="gameMaster"
                                    value={data.gameMaster}
                                    onChange={(e) => setData({ ...data, gameMaster: e.target.value })}
                                    onBlur={() => setIsEditingGameMaster(false)}
                                    autoFocus
                                    variant="outlined"
                                    fullWidth
                                    sx={{ corsor: 'pointer', marginTop: '0.5rem' }}
                                />
                                : <Typography variant="h2" onClick={() => editable && setIsEditingGameMaster(true)}>
                                    {data.gameMaster || 'Navn'}
                                </Typography>
                            }
                        </Box>
                    </Box>
                    <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                        <Box component={FontAwesomeIcon} icon={faScroll} size="2x" sx={{ color: 'primary.main' }} />
                        <Box>
                            <Typography component="span" sx={{ color: 'primary.main' }}>
                                System
                            </Typography>
                            {isEditingSystem ?
                                <TextField
                                    name="system"
                                    value={data.system}
                                    onChange={(e) => setData({ ...data, system: e.target.value })}
                                    onBlur={() => setIsEditingSystem(false)}
                                    autoFocus
                                    variant="outlined"
                                    fullWidth
                                    sx={{ marginTop: '0.5rem' }}
                                />
                                : <Typography variant="h2" onClick={() => editable && setIsEditingSystem(true)}>
                                    {data.system || 'System'}
                                </Typography>
                            }
                        </Box>
                    </Box>
                </Box>
                <ChipCarousel data={data} setData={setData} editable={editable} handleChange={handleChange} />

                {editable ? null : (
                    <Box
                        sx={{
                            '--slider-interest-width': '100%',
                        }}
                    >
                        <InterestSelector disabled />
                    </Box>
                )}
                {editable ? null : (
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: !prevNavigationId || !nextNavigationId ? '1fr' : '1fr 1fr',
                            placeItems: 'space-between',
                            marginBlockStart: '1rem',
                            position: 'relative',
                            width: '100%',
                            rowGap: '1rem',
                        }}
                    >
                        <NavigatePreviousLink previousNavigationId={prevNavigationId} />
                        <NavigateNextLink nextNavigationId={nextNavigationId} />
                        <GoToEventAdministrationButton parentEventId={data.parentEventId} isAdmin={isAdmin} />
                    </Box>
                )}
                {isEditingShortDescription ?
                    <TextField
                        name="shortDescription"
                        value={data.shortDescription}
                        onChange={(e) => setData({ ...data, shortDescription: e.target.value })}
                        onBlur={() => setIsEditingShortDescription(false)}
                        autoFocus
                        variant="outlined"
                        fullWidth
                        sx={{ marginBottom: '1rem' }}
                    />
                    : <Typography
                        sx={{ corsor: 'pointer', ...paragraphStyle, marginBottom: '1rem', textAlign: 'center' }}
                        onClick={() => editable && setIsEditingShortDescription(true)}
                    >
                        {data.shortDescription || 'Kort beskrivelse'}
                    </Typography>
                }
                <Box
                    sx={{ cursor: 'pointer' }}
                    onClick={() => editable && editDescription && editDescription(true)}
                >
                    <MuiMarkdown>{data.description || '# Lang beskrivelse'}</MuiMarkdown>
                </Box>
            </Box>
        </Paper>
        : <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100dvh' }}>
            <Typography variant="h1" marginInlineEnd={2}>
                Henter arrangementet
            </Typography>
            <CircularProgress sx={{ fontSize: '18px' }} />
        </Box>;
};

export default MainEvent;
