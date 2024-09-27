import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import BeginnerIcon from '$lib/components/icons/BeginnerIcon';
import BoardGameIcon from '$lib/components/icons/BoardGameIcon';
import CardGameIcon from '$lib/components/icons/CardGameIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import EnglishIcon from '$lib/components/icons/EnglishIcon';
import LessHoursIcon from '$lib/components/icons/LessHoursIcon';
import MiscGameIcon from '$lib/components/icons/MiscGameIcon';
import MoreHoursIcon from '$lib/components/icons/MoreHoursIcon';
import RoleplayingGameIcon from '$lib/components/icons/RoleplayingGameIcon';
import { GameType } from '$lib/enums';
import type { IconName, IconOption, IconTypes, PoolEvent } from '$lib/types';
import type { Palette } from '@mui/material';
export type ColorProp = keyof Pick<Palette, 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'>;
/**
 * Creates an array of icon options based on various attributes such as adult-only, child-friendly,
 * beginner-friendly, duration, and language support. All options default to `true`.
 * To exclude an option, explicitly set it to `false`.
 *
 * @param {boolean} [adultsOnly=true] - Indicates if the event is for adults only (18+). Defaults to `true`.
 * @param {boolean} [childFriendly=true] - Indicates if the event is child-friendly. Defaults to `true`.
 * @param {boolean} [beginnerFriendly=true] - Indicates if the event is beginner-friendly. Defaults to `true`.
 * @param {boolean} [lessThanThreeHours=true] - Indicates if the event lasts less than three hours. Defaults to `true`.
 * @param {boolean} [moreThanSixHours=true] - Indicates if the event lasts more than six hours. Defaults to `true`.
 * @param {boolean} [possiblyEnglish=true] - Indicates if the event can be run in English. Defaults to `true`.
 *
 * @returns {IconOption[]} An array of icon options with corresponding labels and icon types.
 *
 * @typedef {Object} IconOption
 * @property {string} label - The label describing the icon option.
 * @property {string} icon - The icon type (e.g., 'adultsOnly', 'childFriendly').
 */
export const createIconOptions = (
    adultsOnly = true,
    childFriendly = true,
    beginnerFriendly = true,
    lessThanThreeHours = true,
    moreThanSixHours = true,
    possiblyEnglish = true,
    cardGame = true,
    boardGame = true,
    rolePlaying = true,
    other = true
) => {
    let icons: IconOption[] = [];
    if (childFriendly) icons = [...icons, { label: 'Barnevennlig', iconName: 'childFriendly' }];
    if (possiblyEnglish) icons = [...icons, { label: 'Can be run in English', iconName: 'possiblyEnglish' }];
    if (adultsOnly) icons = [...icons, { label: 'Kun for voksne (18+)', iconName: 'adultsOnly' }];
    if (lessThanThreeHours) icons = [...icons, { label: 'Mindre enn tre timer', iconName: 'lessThanThreeHours' }];
    if (moreThanSixHours) icons = [...icons, { label: 'Mer enn seks timer', iconName: 'moreThanSixHours' }];
    if (beginnerFriendly) icons = [...icons, { label: 'Nybegynnervennlig', iconName: 'beginnerFriendly' }];
    if (cardGame) icons = [...icons, { label: 'Kortspel', iconName: 'cardGame' }];
    if (boardGame) icons = [...icons, { label: 'Brettspel', iconName: 'boardGame' }];
    if (rolePlaying) icons = [...icons, { label: 'Rollespel', iconName: 'rolePlaying' }];
    if (other) icons = [...icons, { label: 'Annet', iconName: 'other' }];
    return icons;
};

/**
 * Converts an icon name (string) to its corresponding React component.
 * This function maps predefined icon strings to their respective icon components.
 *
 * @param {IconName} iconString - The name of the icon as a string. Must be one of the following:
 *  'childFriendly', 'possiblyEnglish', 'adultsOnly', 'lessThanThreeHours', 'moreThanSixHours', 'beginnerFriendly'.
 * @returns {JSX.Element | undefined} - The corresponding icon component if a match is found, otherwise `undefined`.
 *
 * @example
 * const chipOptions = createIconOptions()
 * chipOptions.map((option) => <Chip label={option.label} icon={createIconFromString(option.icon)} />)
 *
 * @example
 * const iconComponent = createIconFromString('childFriendly');
 * // Returns <ChildFriendlyIcon />
 *
 */
export const createIconFromString = (iconString: IconName, color: ColorProp = 'primary') => {
    switch (iconString) {
        case 'childFriendly':
            return <ChildFriendlyIcon color={color} />;
        case 'possiblyEnglish':
            return <EnglishIcon color={color} />;
        case 'adultsOnly':
            return <AdultsOnlyIcon color={color} />;
        case 'lessThanThreeHours':
            return <LessHoursIcon color={color} />;
        case 'moreThanSixHours':
            return <MoreHoursIcon color={color} />;
        case 'beginnerFriendly':
            return <BeginnerIcon color={color} />;
        case 'cardGame':
            return <CardGameIcon color={color} />;
        case 'boardGame':
            return <BoardGameIcon color={color} />;
        case 'rolePlaying':
            return <RoleplayingGameIcon color={color} />;
        case 'other':
            return <MiscGameIcon color={color} />;
        default:
            return undefined;
    }
};
