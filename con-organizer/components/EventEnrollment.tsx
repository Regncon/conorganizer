import { Box, Checkbox } from "@mui/material";
import { EnrollmentOptions } from "@/models/enums";
import { EnrollmentChoice } from "@/models/types";



type Props = {
    enrollmentChoice: EnrollmentChoice;
    handleChoiceChange: (event: React.ChangeEvent<HTMLInputElement>, enrollmentChoice: EnrollmentChoice) => void;
};

const EventEnrollment = ({ enrollmentChoice, handleChoiceChange }: Props) => {
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'row',
                p: 2,
                m: 1,
                borderRadius: 1,
                border: '1px solid #ccc',
                width: '100%',
            }}
            key={enrollmentChoice.id}
        >
            <span>{enrollmentChoice.name}</span>
            <Checkbox
                checked={enrollmentChoice.isEnrolled}
                onChange={(event) => handleChoiceChange(event, enrollmentChoice)}
            />
            {enrollmentChoice.hasGotFirstChoice && <span>Har fått førstevalg</span>}
        </Box>
    );
};

export default EventEnrollment;