import './MainPage.css'
import Navbar from '../../widgets/Navbar/Navbar'
import Card from '../../widgets/Card/Card'
import { useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react'
import Breadcrumb from 'react-bootstrap/Breadcrumb';
import testData from '../../data';

interface Data {
    requestID: number;
    roads: {
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
      }[];

  }

  const MainPage: React.FC = () => {
    const navigate = useNavigate();
    const [data, setData] = useState<Data | null>({ requestID: 0, roads: [] });
    const [minLength, setMinLength] = useState<number | null>(null);
    const fetchData = async (minLength?: string) => {
        try {
            const url = minLength ? `/api/roads?minLength=${minLength}` : '/api/roads';
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`Ошибка при выполнении запроса: ${response.statusText}`);
            }

            const result = await response.json();
            console.log(result); // Проверьте, что данные приходят корректно
            setData(result);
        } catch (error) {
            console.log(testData)
            const result = { ...testData }; // Создаем копию оригинальных данных
            if (minLength) {
                result.roads = testData.roads.filter((roads) => roads.Endofsection-roads.Startofsection >= parseInt(minLength));
            }
            setData(result)
            console.error('ошибка при выполннении запроса:', error);
        }
    };

    const handleMinLengthChange = (value: string) => {
        setMinLength(value !== '' ? parseInt(value) : null);

        // Обновляем URL с использованием navigate
        const minLengthString = value !== '' ? parseInt(value).toString() : '';
        navigate(`?minLength=${minLengthString}`, { replace: true });

        fetchData(minLengthString); // Вызывайте fetchData при изменении maxPrice
    };

    useEffect(() => {
        // Получаем значение minLength  из URL при монтировании компонента
        const urlSearchParams = new URLSearchParams(window.location.search);
        const minLengthParam = urlSearchParams.get('minLength') || '';
        const parsedMinLength = minLengthParam!== null ? parseInt(minLengthParam) : null;
        if (parsedMinLength !== minLength) {
            setMinLength(parsedMinLength);
            fetchData(minLengthParam);
        }
    }, [minLength]);

    return (
        <div>
            <Navbar onMinLengthChange={handleMinLengthChange}/>
            <div className="container">
            <Breadcrumb>
                    <Breadcrumb.Item href="/" active>Главная</Breadcrumb.Item>
                </Breadcrumb>
                <div className="row">
                    {data?.roads?.map((item) => (
                        <div key={item.Roadid} className="col-lg-4 col-md-4 col-sm-12">
                            <Card
                                id={item.Roadid}
                                name={item.Name}
                                image={item.Image}
                                price={item.Price}
                            />
                        </div>
                    ))}
                </div>
            </div>
        </div>
    )

}

export default MainPage 