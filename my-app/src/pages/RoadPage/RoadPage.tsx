import Button from 'react-bootstrap/Button';
import Navbar from '../../widgets/Navbar/Navbar';
import { useParams } from 'react-router-dom';
import { Container, Row, Col } from 'react-bootstrap';
import { useState, useEffect } from 'react';
import Breadcrumb from 'react-bootstrap/Breadcrumb';
import testData from '../../data';

interface RoadData {
    Roadid: number;
    Name: string;
    Trustmanagment: number;
    Length: number;
    Paidlength: number;
    Category: string;
    Numberofstripes: string;
    Speed: number;
    Price: number;
    Image: string;
    Statusroad: string;
    Startofsection: number;
    Endofsection: number;
  }
  const RoadPage: React.FC = () => {
    const { id } = useParams();
    console.log(id)

    const [data, setData] = useState<RoadData | null>(null);

    useEffect(() => {
      // Выполняем запрос при монтировании компонента
      fetchData();
    }, []);


    const fetchData = async () => {
        try {
          const response = await fetch(`/api/roads/${id}`);
          if (!response.ok) {
            throw new Error(`Ошибка при выполнении запроса: ${response.statusText}`);
          }
    
          const result = await response.json();
          setData(result);
        } catch (error) {
            setData(testData.roads[parseInt(id || '0', 10)-1])
          console.error('ошибка при выполнении запроса:', error);
        }
      };
      console.log(data);

    return (
        <div>
            <Navbar />
            <div className="container">
            <Breadcrumb>
                    <Breadcrumb.Item href="/">Главная</Breadcrumb.Item>
                    <Breadcrumb.Item href="#" active>
                        {data?.Name}
                    </Breadcrumb.Item>
                </Breadcrumb>
            <Container style={{ maxWidth: '800px', margin: '20px auto', textAlign: 'center' }}>
                <h1 style={{ fontSize: '2.5rem', margin: '0', fontWeight: 'bold', marginBottom: '20px' }}>{data?.Name}</h1>

                <Row>
                    <Col md={6}>
                        <img
                            src={data?.Image}
                            className="card-img-selected"
                            alt={data?.Name}
                            style={{
                                width: '120%',
                                height: 'auto',
                                display: 'block',
                                border: '2px solid #fff',
                                boxShadow: '0 0 10px rgba(0, 0, 0, 0.3)',
                                borderRadius: '8px 8px 8px 8px',
                                marginBottom: '8px',
                            }}
                        />
                    </Col>
                    <Col md={6} style={{ textAlign: 'left' }}>
                        <ul style={{ listStyleType: 'none', padding: '0', textAlign: 'left', fontSize: '1.1rem' }}>
                            <li>В доверительном управлении: {data?.Trustmanagment} км</li>
                            <li>Начало участка: {data?.Startofsection} км</li>
                            <li>Конец участка: {data?.Endofsection} км</li>
                            <li>Протяженность трассы: {data?.Length} км</li>
                            <li>Протяженность платных участков: {data?.Paidlength} км</li>
                            <li>Категория дороги: {data?.Category}</li>
                            <li>Число полос движения: {data?.Numberofstripes}</li>
                            <li>Разрешенная скорость: до {data?.Speed} км/ч</li>
                            <li>Стоимость проезда: {data?.Price} руб</li>
                            <li style={{ textAlign: 'right' }}>
                                <Button variant="primary">Добавить</Button>
                            </li>
                        </ul>
                    </Col>
                </Row>
            </Container>
            </div>
        </div>
    );
};

export default RoadPage;
