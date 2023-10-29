import { useEffect, useMemo, useState } from 'react';
import { faChessKing, faChild, faDiceD20, faPalette, faUserPlus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { FilterAlt } from '@mui/icons-material';
import { Box, Chip, Typography } from '@mui/material';
import { useSearchParams } from 'next/navigation';
import { useAllEvents } from '@/lib/hooks/UseAllEvents';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { CustomEventTypeNames, GameType, Pool } from '@/models/enums';

type Props = {
    displayPool: Pool;
};

const Filters = ({ displayPool }: Props) => {
    const { events } = useAllEvents();
    console.log(displayPool);
    const [showFilters, setShowFilters] = useState(false);

    const [childFriendly, setChildFriendly] = useState<boolean>(false);
    const [possiblyEnglish, setPossiblyEnglish] = useState<boolean>(false);
    const [volunteersPossible, setVolunteersPossible] = useState<boolean>(false);
    const [beginnerFriendly, setBeginnerFriendly] = useState<boolean>(false);
    const [roleplaying, setRoleplaying] = useState<boolean>(false);
    const [boardgame, setBoardgame] = useState<boolean>(false);
    const [other, setOther] = useState<boolean>(false);

    const filteredEvents = useMemo(() => {
        return events
            ?.filter((conEvent) => conEvent.pool === displayPool)
            .filter((conEvent) => !roleplaying || conEvent.gameType === GameType.roleplaying)
            .filter((conEvent) => !boardgame || conEvent.gameType === GameType.boardgame)
            .filter((conEvent) => !other || conEvent.gameType === GameType.other)
            .filter((conEvent) => !childFriendly || conEvent.childFriendly)
            .filter((conEvent) => !possiblyEnglish || conEvent.possiblyEnglish)
            .filter((conEvent) => !volunteersPossible || conEvent.volunteersPossible)
            .filter((conEvent) => !beginnerFriendly || conEvent.beginnerFriendly);
    }, [
        beginnerFriendly,
        boardgame,
        childFriendly,
        displayPool,
        events,
        other,
        possiblyEnglish,
        roleplaying,
        volunteersPossible,
    ]);

    const handleChanges = new CustomEvent(CustomEventTypeNames.FilterChanges, {
        detail: {
            filteredEvents,
        },
        bubbles: true,
    });
    useEffect(() => {
        document.dispatchEvent(handleChanges);
    }, [filteredEvents]);

    return (
        <>
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
                <Chip
                    icon={<FontAwesomeIcon icon={faDiceD20} color="orangered" />}
                    label="Rollespill"
                    variant={roleplaying ? 'filled' : 'outlined'}
                    onClick={() => {
                        setRoleplaying(!roleplaying);
                        setBoardgame(false);
                        setOther(false);
                    }}
                />
                <Chip
                    icon={<FontAwesomeIcon icon={faChessKing} color="orangered" />}
                    label="Brettspill"
                    variant={boardgame ? 'filled' : 'outlined'}
                    onClick={() => {
                        setBoardgame(!boardgame);
                        setRoleplaying(false);
                        setOther(false);
                    }}
                />
                <Chip
                    label="Annet"
                    icon={<FontAwesomeIcon icon={faPalette} color="orangered" />}
                    variant={other ? 'filled' : 'outlined'}
                    onClick={() => {
                        setOther(!other);
                        setBoardgame(false);
                        setRoleplaying(false);
                    }}
                />
                <Chip
                    icon={<FilterAlt />}
                    label="Andre filtre"
                    variant={showFilters ? 'filled' : 'outlined'}
                    onClick={() => {
                        setShowFilters(!showFilters);
                    }}
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
                        onClick={() => {
                            setChildFriendly(!childFriendly);
                        }}
                    />
                )}
                <Chip
                    icon={<Typography>&nbsp;🇬🇧</Typography>}
                    label="Can be run in English"
                    variant={possiblyEnglish ? 'filled' : 'outlined'}
                    onClick={() => {
                        setPossiblyEnglish(!possiblyEnglish);
                    }}
                />
                <Chip
                    icon={
                        <p>
                            &nbsp;
                            <FontAwesomeIcon icon={faUserPlus} fontSize=".8em" color="#55cc99" />
                        </p>
                    }
                    label="Kan arrangeres av frivillige"
                    variant={volunteersPossible ? 'filled' : 'outlined'}
                    onClick={() => {
                        setVolunteersPossible(!volunteersPossible);
                    }}
                />
                <Chip
                    icon={<Typography>&nbsp;😅</Typography>}
                    label="Nybegynnervennlig"
                    variant={beginnerFriendly ? 'filled' : 'outlined'}
                    onClick={() => {
                        setBeginnerFriendly(!beginnerFriendly);
                    }}
                />
            </Box>
        </>
    );
};

export default Filters;
