'use client';

import { useEffect, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { faChessKing, faChild, faDiceD20, faPalette, faUserPlus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { AccountCircle, FilterAlt } from '@mui/icons-material';
import { useSearchParams } from 'next/navigation';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { GameType, Pool } from '@/models/enums';
import { Box, Chip, Typography } from '../lib/mui';
import EventCardBoundary from './ErrorBoundaries/EventCardBoundary';
import { useAuth } from './AuthProvider';
import EventCard from './EventCard';
import PoolSelector from './PoolSelector';

const EventList = () => {
    const user = useAuth();
    const { events, loading } = useAllEvents();
    const [displayPool, setDisplayPool] = useState<Pool>(Pool.FridayEvening);
    const [showFilters, setShowFilters] = useState(false);
    const [showUnpublished, setShowUnpublished] = useState(false);
    const { userSettings } = useUserSettings(user?.uid);
    const searchParams = useSearchParams();

    // const [mySignups, setMySignups] = useState<boolean>(false);
    const [childFriendly, setChildFriendly] = useState<boolean>(false);
    const [possiblyEnglish, setPossiblyEnglish] = useState<boolean>(false);
    // const [adultsOnly, setAdultsOnly] = useState<boolean>(false);
    const [volunteersPossible, setVolunteersPossible] = useState<boolean>(false);
    const [beginnerFriendly, setBeginnerFriendly] = useState<boolean>(false);
    const [roleplaying, setRoleplaying] = useState<boolean>(false);
    const [boardgame, setBoardgame] = useState<boolean>(false);
    const [other, setOther] = useState<boolean>(false);

    useEffect(() => {
        setShowUnpublished(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

    const search = searchParams.get('pool') as keyof typeof Pool;
    useEffect(() => {
        if (search) {
            setDisplayPool(Pool[search]);
        }
    }, [search]);

    return (
        <>
            <PoolSelector handlePoolChange={(pool) => setDisplayPool(pool)} />

            <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20">
                {loading ? <Typography variant="body1">Loading...</Typography> : null}
                <Box
                    sx={{
                        display: 'flex',
                        gap: '.5em',
                        flexGrow: '1',
                        justifyContent: 'center',
                        width: '100%',
                        flexWrap: 'wrap',
                    }}
                >
                    {/* <Chip
                        label="Alle"
                        variant={
                            !(
                                roleplaying ||
                                boardgame ||
                                other ||
                                childFriendly ||
                                beginnerFriendly ||
                                possiblyEnglish ||
                                // adultsOnly ||
                                volunteersPossible
                            )
                                ? 'filled'
                                : 'outlined'
                        }
                    /> */}
                    {/* <Chip label="Mine p&aring;meldinger" variant="outlined" icon={<AccountCircle />} /> */}
                    <Chip
                        icon={<FontAwesomeIcon icon={faDiceD20} color="orangered" />}
                        label="Rollespill"
                        variant={roleplaying ? 'filled' : 'outlined'}
                        onClick={() => setRoleplaying(!roleplaying)}
                    />
                    <Chip
                        icon={<FontAwesomeIcon icon={faChessKing} color="orangered" />}
                        label="Brettspill"
                        variant={boardgame ? 'filled' : 'outlined'}
                        onClick={() => setBoardgame(!boardgame)}
                    />
                    <Chip
                        label="Annet"
                        icon={<FontAwesomeIcon icon={faPalette} color="orangered" />}
                        variant={other ? 'filled' : 'outlined'}
                        onClick={() => setOther(!other)}
                    />
                    <Chip
                        icon={<FilterAlt />}
                        label="Andre filtre"
                        variant={showFilters ? 'filled' : 'outlined'}
                        onClick={() => setShowFilters(!showFilters)}
                    />
                </Box>
                <Box
                    display={showFilters ? 'flex' : 'none'}
                    sx={{ gap: '.5em', flexGrow: '1', justifyContent: 'center', width: '100%', flexWrap: 'wrap' }}
                >
                    {displayPool === 'Søndag Morgen' && (
                        <Chip
                            icon={
                                <p>
                                    &nbsp;
                                    <FontAwesomeIcon icon={faChild} fontSize="1em" color="#00aaff" />
                                </p>
                            }
                            label="Barnevennlig"
                            variant={childFriendly ? 'filled' : 'outlined'}
                            onClick={() => setChildFriendly(!childFriendly)}
                        />
                    )}
                    <Chip
                        icon={<Typography>&nbsp;🇬🇧</Typography>}
                        label="Can be run in English"
                        variant={possiblyEnglish ? 'filled' : 'outlined'}
                        onClick={() => setPossiblyEnglish(!possiblyEnglish)}
                    />
                    {/* <Chip
                        icon={<Typography fontSize=".9rem">&nbsp;🔞</Typography>}
                        label="Kun for voksne"
                        variant={adultsOnly ? 'filled' : 'outlined'}
                        onClick={() => setAdultsOnly(!adultsOnly)}
                    /> */}
                    <Chip
                        icon={
                            <p>
                                &nbsp;
                                <FontAwesomeIcon icon={faUserPlus} fontSize=".8em" color="#55cc99" />
                            </p>
                        }
                        label="Kan arrangeres av frivillige"
                        variant={volunteersPossible ? 'filled' : 'outlined'}
                        onClick={() => setVolunteersPossible(!volunteersPossible)}
                    />
                    <Chip
                        icon={<Typography>&nbsp;😅</Typography>}
                        label="Nybegynnervennlig"
                        variant={beginnerFriendly ? 'filled' : 'outlined'}
                        onClick={() => setBeginnerFriendly(!beginnerFriendly)}
                    />
                </Box>
                {events
                    ?.filter((conEvent) => conEvent.pool === displayPool)
                    // .filter((conEvent) => !mySignups || conEvent.published)
                    .filter((conEvent) => !roleplaying || conEvent.gameType === GameType.roleplaying)
                    .filter((conEvent) => !boardgame || conEvent.gameType === GameType.boardgame)
                    .filter((conEvent) => !other || conEvent.gameType === GameType.other)
                    .filter((conEvent) => !childFriendly || conEvent.childFriendly)
                    .filter((conEvent) => !possiblyEnglish || conEvent.possiblyEnglish)
                    // .filter((conEvent) => !adultsOnly || conEvent.adultsOnly)
                    .filter((conEvent) => !volunteersPossible || conEvent.volunteersPossible)
                    .filter((conEvent) => !beginnerFriendly || conEvent.beginnerFriendly)
                    .filter((conEvent) => showUnpublished || conEvent.published)
                    .toSorted((a, b) => a.sortingIndex - b.sortingIndex)
                    .map((conEvent) => (
                        <ErrorBoundary FallbackComponent={EventCardBoundary} key={conEvent.id}>
                            <EventCard conEvent={conEvent} />
                        </ErrorBoundary>
                    ))}
            </Box>
        </>
    );
};

export default EventList;
