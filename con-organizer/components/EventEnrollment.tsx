import WarningIcon from '@mui/icons-material/Warning';
import { Box, Checkbox, TableCell, TableRow, Tooltip } from '@mui/material';
import { EnrollmentChoice } from '@/models/types';

type Props = {
    enrollmentChoice: EnrollmentChoice;
    handleChoiceChange: (event: React.ChangeEvent<HTMLInputElement>, enrollmentChoice: EnrollmentChoice) => void;
};

const EventEnrollment = ({ enrollmentChoice, handleChoiceChange }: Props) => {
    return (
        <TableRow key={enrollmentChoice.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
            <TableCell align="right">{enrollmentChoice.name}</TableCell>
            <TableCell align="right"
                sx={{
                    display: 'flex',
                    flexDirection: 'row',
                    alignItems: 'center',
                    minWidth: '200px',
                }}
            >
                <span>Fått plass:</span>
                
                <Checkbox
                    checked={enrollmentChoice.isEnrolled}
                    onChange={(event) => handleChoiceChange(event, enrollmentChoice)}
                />
            </TableCell>
            <TableCell align="right">
                {enrollmentChoice.ticketType && <span>{enrollmentChoice.ticketType}</span>}
            </TableCell>
            <TableCell align="right">
                {enrollmentChoice.hasGotFirstChoice && (
                    <Tooltip title="Har fått førstevalg">
                        <WarningIcon sx={{ color: 'warning.main', mr: 1 }} />
                    </Tooltip>
                )}
            </TableCell>

            <span></span>
        </TableRow>
    );
};

export default EventEnrollment;
