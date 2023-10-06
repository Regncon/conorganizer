'use client';

import { Box, Card, CardContent, Checkbox, Divider } from '@mui/material';
import { useAllEnrollmentChoices } from '@/lib/hooks/UseAllEnrollmentChoices';
import { EnrollmentOptions } from '@/models/enums';

type Props = { id: string };

const EventEnrolementList = ({ id }: Props) => {
    const { enrollmentChoices, loadingEnrollmentChoices } = useAllEnrollmentChoices(id || '');

    console.log(enrollmentChoices);

    return (
        <Card>
            <CardContent>
                {enrollmentChoices
                    ?.toSorted((a, b) => a.choice - b.choice)
                    .map((enrollmentChoice) => (
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
                            <span>{EnrollmentOptions[enrollmentChoice.choice]}</span>
                            <span>{enrollmentChoice.name}</span>
                            <Checkbox value={enrollmentChoice.isEnrolled} />
                        </Box>
                    ))}
            </CardContent>
        </Card>
    );
};

export default EventEnrolementList;
