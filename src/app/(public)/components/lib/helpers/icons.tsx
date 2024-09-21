import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import BeginnerIcon from '$lib/components/icons/BeginnerIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import EnglishIcon from '$lib/components/icons/EnglishIcon';
import LessHoursIcon from '$lib/components/icons/LessHoursIcon';
import MoreHoursIcon from '$lib/components/icons/MoreHoursIcon';
import type { IconName, IconOption, IconTypes, PoolEvent } from '$lib/types';

export const createIconArray = ({
    adultsOnly,
    childFriendly,
    beginnerFriendly,
    lessThanThreeHours,
    moreThanSixHours,
    possiblyEnglish,
}: IconTypes) => {
    let icons: IconOption[] = [];
    if (childFriendly) icons = [...icons, { label: 'Barnevennlig', icon: 'childFriendly' }];
    if (possiblyEnglish) icons = [...icons, { label: 'Can be run in English', icon: 'possiblyEnglish' }];
    if (adultsOnly) icons = [...icons, { label: 'Kun for voksne (18+)', icon: 'adultsOnly' }];
    if (lessThanThreeHours) icons = [...icons, { label: 'Mindre enn tre timer', icon: 'lessThanThreeHours' }];
    if (moreThanSixHours) icons = [...icons, { label: 'Mer enn seks timer', icon: 'moreThanSixHours' }];
    if (beginnerFriendly) icons = [...icons, { label: 'Nybegynnervennlig', icon: 'beginnerFriendly' }];
    return icons;
};

export const createIconColor = (iconString: IconName) => { };
export const createIconFromString = (iconString: IconName) => {
    switch (iconString) {
        case 'childFriendly':
            return <ChildFriendlyIcon />;
        case 'possiblyEnglish':
            return <EnglishIcon />;
        case 'adultsOnly':
            return <AdultsOnlyIcon />;
        case 'lessThanThreeHours':
            return <LessHoursIcon />;
        case 'moreThanSixHours':
            return <MoreHoursIcon />;
        case 'beginnerFriendly':
            return <BeginnerIcon />;
        default:
            return undefined;
    }
};
