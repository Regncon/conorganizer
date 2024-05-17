'use client';
import { faChevronLeft } from '@fortawesome/free-solid-svg-icons/faChevronLeft';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import IconButton from '@mui/material/IconButton';

const BackButton = () => {
    return (
        <IconButton
            onClick={() => {
                history.back();
            }}
        >
            <FontAwesomeIcon icon={faChevronLeft} fixedWidth />
        </IconButton>
    );
};

export default BackButton;
