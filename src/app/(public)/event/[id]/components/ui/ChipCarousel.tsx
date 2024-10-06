import { createIconOptions, createIconFromString } from '$app/(public)/components/lib/helpers/icons';
import { GameType } from '$lib/enums';
import type { ConEvent, IconName, PoolEvent } from '$lib/types';
import {
    Box,
    Button,
    Chip,
    debounce,
    FormControl,
    FormControlLabel,
    Radio,
    RadioGroup,
    TextField,
} from '@mui/material';
import { useCallback, useEffect, useRef, useState, type FormEventHandler, type PropsWithChildren } from 'react';

type SelectedTag = {
    selected: boolean;
    label: string;
    iconName: IconName;
};

type Props = {
    data: PoolEvent;
    setData: (data: PoolEvent) => void;
    editable: boolean;
    handleChange?: (data: Partial<ConEvent>) => Promise<void>;
};
const ChipCarousel = ({ data, editable, setData, handleChange }: Props) => {
    const chipIcons = createIconOptions(
        data.adultsOnly,
        data.childFriendly,
        data.beginnerFriendly,
        data.lessThanThreeHours,
        data.moreThanSixHours,
        data.possiblyEnglish,
        data.gameType
    );
    const [isEditingTags, setIsEditingTags] = useState<boolean>(false);
    const [selectedTags, setSelectedTags] = useState<SelectedTag[]>(
        chipIcons.map((chipIcon) => {
            return { ...chipIcon, selected: data.icons?.some((icon) => icon.iconName === chipIcon.iconName) ?? false };
        })
    );

    const originalIconLength = data?.icons?.length ?? 0; // Get the length of the original icons

    const handleBoxAction = () => {
        editable && setIsEditingTags(!isEditingTags);
    };
    const updateConEvents = useCallback(
        debounce((selectedTagsToConEventData: Partial<ConEvent>) => {
            console.log('selectedTagsToConEventData: ', selectedTagsToConEventData);

            handleChange?.(selectedTagsToConEventData);
        }, 1000),
        []
    );
    const handleChipClick = (clickedChips: SelectedTag) => {
        const selectedTagsToConEventData: Partial<ConEvent> = {
            ...selectedTags.reduce((acc, value) => {
                return { ...acc, [value.iconName]: value.selected };
            }, {}),
            [clickedChips.iconName]: clickedChips.selected,
        };

        updateConEvents(selectedTagsToConEventData);

        setSelectedTags((prev) => {
            return [
                ...prev.map((chip) =>
                    chip.iconName === clickedChips.iconName ? { ...chip, selected: clickedChips.selected } : chip
                ),
            ];
        });
    };
    const handleChangeGameType: FormEventHandler = (e) => {
        const target = e.target as HTMLInputElement;
        const radioName = target.name as GameType;
        const radioTextContent = target.labels?.[0].textContent;
        if (radioTextContent) {
            const deletedGameTypeSelectedTags = selectedTags.filter(
                (selectedTag) => Object.values(GameType).includes(selectedTag.iconName as GameType) === false
            );

            setSelectedTags([
                ...deletedGameTypeSelectedTags,
                {
                    iconName: radioName,
                    label: radioTextContent,
                    selected: true,
                },
            ]);
        }
    };

    return isEditingTags ?
            <Box sx={{ display: 'grid' }}>
                <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
                    {selectedTags.map((iconOption, index) => {
                        return (
                            <Chip
                                onClick={() => {
                                    handleChipClick({ ...iconOption, selected: !iconOption.selected });
                                }}
                                label={iconOption.label}
                                key={`${iconOption.label}-${index}`}
                                color={iconOption.selected ? 'secondary' : 'primary'}
                                variant="outlined"
                                icon={createIconFromString(iconOption.iconName)}
                                disabled={Object.values(GameType).includes(iconOption.iconName as GameType)}
                            />
                        );
                    })}
                </Box>
                <FormControl
                    sx={{ placeItems: 'center' }}
                    onChange={(e) => {
                        console.log(e, 'e');
                        handleChangeGameType(e);
                    }}
                >
                    <RadioGroup
                        value={
                            selectedTags.find((selectedTab) =>
                                Object.values(GameType).includes(selectedTab.iconName as GameType)
                            )?.iconName
                        }
                        aria-labelledby="demo-controlled-radio-buttons-group"
                        name="controlled-radio-buttons-group"
                    >
                        <FormControlLabel
                            value="rolePlaying"
                            control={<Radio name="rolePlaying" />}
                            label="rollespel"
                        />
                        <FormControlLabel value="boardGame" control={<Radio name="boardGame" />} label="Brettspel" />
                        <FormControlLabel value="cardGame" control={<Radio name="cardGame" />} label="Kortspel" />
                        <FormControlLabel value="other" control={<Radio name="other" />} label="Annet" />
                    </RadioGroup>
                </FormControl>
                <Button
                    onClick={() => {
                        setIsEditingTags(false);
                    }}
                >
                    Tilbake til karusell
                </Button>
            </Box>
        :   <Box
                sx={{
                    display: 'flex',
                    gap: '.5em',
                    overflow: 'hidden',
                    width: '100%',
                    position: 'relative',
                }}
                onClick={handleBoxAction}
                onBlur={handleBoxAction}
            >
                <Box
                    sx={{
                        display: 'flex',
                        gap: '.5em',
                        animation: 'scrollX 7s linear infinite',
                        '@keyframes scrollX': {
                            '0%': {
                                transform: 'translateX(0)',
                            },
                            '100%': {
                                transform: 'translateX(-50%)',
                            },
                        },
                    }}
                >
                    {[...(data?.icons ?? []), ...(data?.icons ?? [])].map((iconOption, index) => (
                        <Chip
                            label={iconOption.label}
                            key={`${iconOption.label}-${index}`}
                            color="primary"
                            variant="outlined"
                            icon={createIconFromString(iconOption.iconName)}
                            sx={{
                                marginRight: index === originalIconLength - 1 ? '2rem' : 'unset',
                                marginLeft: index === 0 ? '2.5rem' : 'unset',
                            }}
                        />
                    ))}
                </Box>
            </Box>;
};

export default ChipCarousel;
