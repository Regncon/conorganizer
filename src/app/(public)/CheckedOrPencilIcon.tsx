'use client';

import { faCircleCheck } from '@fortawesome/free-solid-svg-icons/faCircleCheck';
import { faPencil } from '@fortawesome/free-solid-svg-icons/faPencil';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Box from '@mui/material/Box';

type Props = {
    myEventBarSubmitted: boolean;
};

const CheckedOrPencilIcon = ({ myEventBarSubmitted }: Props) => {
    const circleCheckOrPencilIcon = myEventBarSubmitted ? faCircleCheck : faPencil;
    const SuccessOrWarningColor = myEventBarSubmitted ? 'success.main' : 'warning.main';
    return (
        <Box
            component={FontAwesomeIcon}
            icon={circleCheckOrPencilIcon}
            sx={{ color: SuccessOrWarningColor }}
            size="2x"
        />
    );
};

export default CheckedOrPencilIcon;
