import { createIconOptions, createIconFromString } from '$app/(public)/components/lib/helpers/icons';
import type { ConEvent, IconName, PoolEvent } from '$lib/types';
import { Box, Chip, debounce, TextField } from '@mui/material';
import { useCallback, useEffect, useRef, useState, type PropsWithChildren } from 'react';

type SelectedTags = {
    selected: boolean;
    label: string;
    iconName: IconName;
};

const chipIcons = createIconOptions();
type Props = {
    data: PoolEvent;
    setData: (data: PoolEvent) => void;
    editable: boolean;
    handleChange?: (data: Partial<ConEvent>) => Promise<void>;
};
const ChipCarousel = ({ data, editable, setData, handleChange }: Props) => {
    const [isEditingTags, setIsEditingTags] = useState<boolean>(false);
    const [selectedTags, setSelectedTags] = useState<SelectedTags[]>(
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
    const handleChipClick = (clickedChips: SelectedTags) => {
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

    return isEditingTags ?
            <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
                {selectedTags.map((iconOption, index) => (
                    <Chip
                        onClick={() => {
                            handleChipClick({ ...iconOption, selected: !iconOption.selected });
                        }}
                        label={iconOption.label}
                        key={`${iconOption.label}-${index}`}
                        color={iconOption.selected ? 'secondary' : 'primary'}
                        variant="outlined"
                        icon={createIconFromString(iconOption.iconName)}
                    />
                ))}
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
