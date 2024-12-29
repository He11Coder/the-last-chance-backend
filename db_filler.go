package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"mainService/internal/domain"
)

type userInfoWrapper struct {
	user              *domain.ApiUserInfo
	userImagePath     string
	userBackImagePath string
}

type petInfoWrapper struct {
	pet        *domain.ApiPetInfo
	ownerID    string
	avatarPath string
}

var serviceURL = "http://localhost:8081"

func addUser(user *domain.ApiUserInfo, userImagePath, userBackImagePath string) (*domain.LoginResponse, error) {
	addURL := serviceURL + "/register"

	userImage, err := os.ReadFile(userImagePath)
	if err != nil {
		return nil, err
	}

	userImageBase64 := base64.StdEncoding.EncodeToString(userImage)

	userBackImage, err := os.ReadFile(userBackImagePath)
	if err != nil {
		return nil, err
	}

	userBackImageBase64 := base64.StdEncoding.EncodeToString(userBackImage)

	user.UserImage = userImageBase64
	user.UserBackImage = userBackImageBase64

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errResp := map[string]interface{}{}
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%v", errResp["message"])
	}

	logResp := new(domain.LoginResponse)
	err = json.Unmarshal(body, logResp)
	if err != nil {
		return nil, err
	}

	return logResp, nil
}

func fillDatabaseWithUsers() error {
	usersToAdd := []userInfoWrapper{}

	user := &domain.ApiUserInfo{
		Login:    "happy_man",
		Password: "qwerty123",
		Username: "Дмитрий Сергеев",
		Contacts: "TG: @DimonSerg",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/1.jpg", userBackImagePath: "assets/to_fill/headers/1.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "serious_grandpa",
		Password: "qwerty123",
		Username: "Василий Александрович Плотников",
		Contacts: "Тел.: +79134567989",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/2.png", userBackImagePath: "assets/to_fill/headers/1.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "cool_raccoon",
		Password: "qwerty123",
		Username: "Андрей",
		Contacts: "VK: /cool_raccoon",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/3.jpg", userBackImagePath: "assets/to_fill/headers/2.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "lovely_girl",
		Password: "qwerty123",
		Username: "Екатерина Смолина",
		Contacts: "OK: /your_queen",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/4.jpg", userBackImagePath: "assets/to_fill/headers/2.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "pro_designer",
		Password: "qwerty123",
		Username: "Василий Смолин",
		Contacts: "WhatsApp: +79830032076",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/5.jpg", userBackImagePath: "assets/to_fill/headers/3.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "kesha_official",
		Password: "qwerty123",
		Username: "Иннокентий Семенович Радулов",
		Contacts: "Тел.: +79851321717",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/6.png", userBackImagePath: "assets/to_fill/headers/3.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "toha_top",
		Password: "qwerty123",
		Username: "Антон Марты",
		Contacts: "Viber: +79051945427",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/7.jpg", userBackImagePath: "assets/to_fill/headers/4.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "real_zigmund",
		Password: "qwerty123",
		Username: "Зигмунд Робертович Пролетарский",
		Contacts: "Тел.: +79111120102",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/8.png", userBackImagePath: "assets/to_fill/headers/4.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "super_jackson",
		Password: "qwerty123",
		Username: "Николас Джексон",
		Contacts: "TG: @SuperJackson",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/9.png", userBackImagePath: "assets/to_fill/headers/5.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "usach",
		Password: "qwerty123",
		Username: "Федор Меркурьев",
		Contacts: "Тел.: +79138971232",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/10.jpg", userBackImagePath: "assets/to_fill/headers/5.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "evil_vector",
		Password: "qwerty123",
		Username: "Виктор Баринов",
		Contacts: "TG: @Vector",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/11.jpg", userBackImagePath: "assets/to_fill/headers/6.jpg"})

	user = &domain.ApiUserInfo{
		Login:    "leha_top",
		Password: "qwerty123",
		Username: "Алексей",
		Contacts: "WhatsApp: +79537551312",
	}
	usersToAdd = append(usersToAdd, userInfoWrapper{user: user, userImagePath: "assets/to_fill/users/12.png", userBackImagePath: "assets/to_fill/headers/6.jpg"})

	for _, u := range usersToAdd {
		resp, err := addUser(u.user, u.userImagePath, u.userBackImagePath)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", *resp)
	}

	return nil
}

func addPet(pet *domain.ApiPetInfo, avatarPath string, userID string) (string, error) {
	addURL := serviceURL + "/add_pet/" + userID

	avatar, err := os.ReadFile(avatarPath)
	if err != nil {
		return "", err
	}

	avatarBase64 := base64.StdEncoding.EncodeToString(avatar)
	pet.PetAvatar = avatarBase64

	jsonPet, err := json.Marshal(pet)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonPet))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		errResp := map[string]interface{}{}
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("%v", errResp["message"])
	}

	addResp := new(domain.ApiPetInfo)
	err = json.Unmarshal(body, addResp)
	if err != nil {
		return "", err
	}

	return addResp.PetID, nil
}

func fillDatabaseWithPets() error {
	petsToAdd := []petInfoWrapper{}

	pet := &domain.ApiPetInfo{
		TypeOfAnimal: "медведь",
		Name:         "Михаил",
		Info:         "Медведь-гризли, возраст 15 лет. Очень любит рыбу (любую, но свежую).\n\nАллергия на консервы. Ласковый, любит, когда его гладят по голове и чешут за ушком. Необходим дневной сон, хотя бы пару часов.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cf1e6eda565556e40fa16", avatarPath: "assets/to_fill/pets/bear.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "олень",
		Name:         "Рудольф",
		Info:         "Благородный олень с красноватым носом. Возраст 7 лет. Игривый. Из еды предпочитает клевер, желуди, каштаны.\n\nАккуратней, у него травмирована правая передняя лапка, наступили случайно. Еще лечимся.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd3feda565556e40fa17", avatarPath: "assets/to_fill/pets/deer.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "кошка",
		Name:         "Шелли",
		Info:         "Абиссинская кошка. Возраст 3 года. Негостеприимная, шипит на незнакомцев, но если подкормить, то добреет.\nНеобходим дневной сон после обеда.\n\nИз еды предпочитает баварские сосиски или кильку в томатном соусе.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd3feda565556e40fa18", avatarPath: "assets/to_fill/pets/cat1.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "собака",
		Name:         "Мухтар",
		Info:         "Немецкая овчарка. Возраст преклонный - 12 лет. Нуждается в заботе и внимании, бережном обращении.\nКормить мясом с низким содержанием жира! Из любимого мяса - индейка и говядина.\nНежная шерсть, любит поласковиться.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd40eda565556e40fa19", avatarPath: "assets/to_fill/pets/dog1.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "белка",
		Name:         "Стрелка",
		Info:         "Белка обыкновенная, подобрал на прогулке в лесу. Точный возраст мне неизвестен, врачи сказали, что около 5 лет. Необходимые прививки стоят.\n\nОчень любит гостей, особенно, если они приносят что-то вкусненькое: еловые шишки, кедровые и лесные орехи, семена пихты.\nАллергия на грибы.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd40eda565556e40fa1a", avatarPath: "assets/to_fill/pets/squirrel.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "змея",
		Name:         "Мила",
		Info:         "Кустарниковая гадюка, возраст 2 года. Очень ласковая, нежная, любит обниматься.\nПочти не кусается, но если кусается, то смертельно. Быть осторожным.\n\nЛюбимые лакомства - лягушки и ящерицы. Не выносит улиток.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd40eda565556e40fa1b", avatarPath: "assets/to_fill/pets/snake.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "собака",
		Name:         "Джесси",
		Info:         "Далматин-девочка. 7 лет. Любит играть на свежем воздухе, необходимо выгуливать хотя бы 2-3 раза в день.\n\nАллергия на баранину и молоко.\n\nИгривая и ласковая девочка, очень гостеприимная и активная. Любит сладкое.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd40eda565556e40fa1c", avatarPath: "assets/to_fill/pets/dog2.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "кошка",
		Name:         "Тиффани",
		Info:         "Кошка породы Жоффруа. Возраст 5 лет. Часто болеет конъюктивитом, обязательно мыть руки перед контактом!\nХарактер непростой, но если с ней подружиться, то она довольно ласковая и милая.\nТакже очень боится воды, купать ее довольно тяжело, рекомендации дам лично. Зато ест все подряд.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd41eda565556e40fa1d", avatarPath: "assets/to_fill/pets/cat2.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "кошка",
		Name:         "Делайла",
		Info:         "Британская короткошерстная кошка, 3 года отроду.\nСтерилизованная.\nБоится незнакомых людей, поэтому нужно будет время, чтобы привыкнуть к новому человеку.\nЕсть проблемы с выпадающей шерстью.\n\nИз еды любит паштет, свежее сырое куриное мясо и молоко. Аллергия на корм из магазина.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676cfd41eda565556e40fa1d", avatarPath: "assets/to_fill/pets/cat3.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "черепаха",
		Name:         "Лео",
		Info:         "Среднеазиатская черепаха. Взрослая - 30 лет. Несмотря на возраст, здоровье очень крепкое.\nВ основном предпочитает сидеть дома в своем аквариуме, но иногда нужно выносить ее на улицу, на свежий воздух.\nПредпочитает рачков, червей или рыбу. Без большого желания, но ест и сухой корм.\nНеобходимо кормить 1 раз в 2-3 дня",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676d016beda565556e40fa1e", avatarPath: "assets/to_fill/pets/turtle.png"})

	pet = &domain.ApiPetInfo{
		TypeOfAnimal: "лиса",
		Name:         "Айгуль",
		Info:         "Маленький лисенок, прибился к дому. Возраст приблизительно 2 года. Все прививки стоят, но есть проблемы с иммунитетом.\nНуждается в бережном уходе и правильном питании. Предпочитает мелких грызунов, например, мышей; жуков. Из деликатесов - мелкие птички, фрукты и различные слакие плоды.",
	}
	petsToAdd = append(petsToAdd, petInfoWrapper{pet: pet, ownerID: "676d0ab84411455f3bffb58e", avatarPath: "assets/to_fill/pets/fox.png"})

	for _, p := range petsToAdd {
		petID, err := addPet(p.pet, p.avatarPath, p.ownerID)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", petID)
	}

	return nil
}

func addService(serv *domain.ApiService) (string, error) {
	addURL := serviceURL + "/add_service/" + serv.UserID

	jsonServ, err := json.Marshal(serv)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonServ))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		errResp := map[string]interface{}{}
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("%v", errResp["message"])
	}

	addResp := new(domain.ApiService)
	err = json.Unmarshal(body, addResp)
	if err != nil {
		return "", err
	}

	return addResp.ServiceID, nil
}

func fillDatabaseWithMasterServices() error {
	servsToAdd := []*domain.ApiService{}

	serv := &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cf1e6eda565556e40fa16",
		Title:       "Выгул медведя",
		Price:       10000,
		Description: "Нужно выгулять медведя в мое отсутствие.\nЦена приблизительная, понимаю, что работа сложная и уникальная, поэтому возможен торг.\n\nРабота постоянная. Каждый понедельник с 12 до 14 часов. Информацию о медведе см. в его профиле.",
		PetIDs:      []string{"676d2599daf6bf4efe109efb"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd3feda565556e40fa17",
		Title:       "Уход за оленем",
		Price:       3000,
		Description: "Требуется уход за моим оленешей Рудольфом. Нужно приходить каждый день, чистить ему лапы, обновлять еду в миске.\n\nЦена указана за один ваш визит.\nВсе подробности лично.",
		PetIDs:      []string{"676d259adaf6bf4efe109efc"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd3feda565556e40fa18",
		Title:       "Стерилизация кошки",
		Price:       2500,
		Description: "Разовая услуга. Нужно стерилизовать кошку. Описание кошки в ее профиле, цена строгая, менять не буду.\n\nМесто и время можем обсудить лично, см. мой контакт в профиле.",
		PetIDs:      []string{"676d259adaf6bf4efe109efd"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd40eda565556e40fa19",
		Title:       "Нужен грумер для собаки",
		Price:       4000,
		Description: "Моя овчарка нуждается в груминге, чистке ушей и стрижке когтей.\nУслугу нужно повторять раз в месяц, подробности можем обсудить лично.\n\nЦена предложена выше рынка, поскольку понимаю, что работать с большой собакой тяжелее.",
		PetIDs:      []string{"676d259adaf6bf4efe109efe"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd40eda565556e40fa1a",
		Title:       "Требуется поставщик орехов",
		Description: "Моя домашняя белка очень любит орехи. Съедает она их очень быстро, поэтому требуется поставлять орехи каждую неделю.\n\nАдрес, объемы поставок, цену и дополнительные условия обсудим лично, см. контакт в профиле.\n\nP.S. Интересуют лесные и кедровые орехи.",
		PetIDs:      []string{"676d259adaf6bf4efe109eff"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd40eda565556e40fa1b",
		Title:       "Требуется чистильщик змеиной коробки",
		Description: "Цена договорная, требуется уборщик за змеей. Работа опасная, поэтому готов предложить хорошие деньги, пишите.\n\nИнформацию о змее можете посмотреть в ее профиле.\nИз обязанностей: чистка коробки, обновление корма в миске (корм дома есть).\n\nВремя обсудим лично.",
		PetIDs:      []string{"676d259bdaf6bf4efe109f00"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd40eda565556e40fa1c",
		Title:       "Сидельщик с собакой на выходные",
		Price:       5000,
		Description: "Уезжаю на выходные из города, нужен собакоситтер. Цена указана за оба выходных дня, работать нужно будет с 12 до 18.\n\nТребуется выгулять собаку дважды: в начале рабочего дня и в конце.\nПосле прогулки необходимо обновить миску с едой и водой, все продукты есть дома, подробности лично.",
		PetIDs:      []string{"676d259bdaf6bf4efe109f01"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676cfd41eda565556e40fa1d",
		Title:       "Сидельщик для кошек по будням",
		Description: "Цену не указываю, договоримся лично. Дома две кошки, читайте про них в их профилях.\n\nНеобходимо приходить в обеденное время каждый будний день, обновлять им корм в мисках, прибираться за ними в квартире, если нужно.\nКорм дома есть, все необходимое для уборки тоже. Детали лично.",
		PetIDs:      []string{"676d259bdaf6bf4efe109f02", "676d259bdaf6bf4efe109f03"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676d016beda565556e40fa1e",
		Title:       "Уборка аквариума с черепахой",
		Price:       1000,
		Description: "Цена за раз. Адрес и время обсудим лично.\n\nНужно будет приходить раз в 2 дня (в вечернее или дневное время) и убираться в аквариуме за черепахой. Работа абсолютно нетрудная, пишите.",
		PetIDs:      []string{"676d259bdaf6bf4efe109f04"},
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Customer,
		UserID:      "676d0ab84411455f3bffb58e",
		Title:       "Личный ветеринар для лисы",
		Price:       5000,
		Description: "Более подробную информацию о лисе смотрите в ее прикрепленном профиле.\nЦена указана за один ваш визит. Возможен торг.\nНужно будет приходить ко мне домой (приблизительно раз в неделю), отслеживать состояние питомца, принимать необходимые меры, выписывать лечение.\n\nПодробности в ЛС.",
		PetIDs:      []string{"676d259bdaf6bf4efe109f05"},
	}
	servsToAdd = append(servsToAdd, serv)

	for _, s := range servsToAdd {
		servID, err := addService(s)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", servID)
	}

	return nil
}

func fillDatabaseWithSlaveServices() error {
	servsToAdd := []*domain.ApiService{}

	serv := &domain.ApiService{
		Type:        domain.Provider,
		UserID:      "676d4ccb20e233c014ef4a00",
		Title:       "Посижу с вашей собакой",
		Price:       1000,
		Description: "Цена указана за час работы\nМогу выгулять собаку на улице, посидеть с ней дома, поиграть, накормить, прибрать за собакой.\n\nP.S. По договоренности могу ухаживать за котами и кошками.",
	}
	servsToAdd = append(servsToAdd, serv)

	serv = &domain.ApiService{
		Type:        domain.Provider,
		UserID:      "676d4ce820e233c014ef4a01",
		Title:       "Ветеринар для кота",
		Description: "Предлагаю ветеринарские услуги для животных, прежде всего для кошек и котов. Цена договорная, зависит от вида услуги.\n\nКастрация/стерилизация\nПрививки\nКонсультации\nОформление документов",
	}
	servsToAdd = append(servsToAdd, serv)

	for _, s := range servsToAdd {
		servID, err := addService(s)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", servID)
	}

	return nil
}

func main() {
	/*err := fillDatabaseWithUsers()
	if err != nil {
		fmt.Println(err)
		return
	}*/

	/*err := fillDatabaseWithPets()
	if err != nil {
		fmt.Println(err)
		return
	}*/

	/*err := fillDatabaseWithMasterServices()
	if err != nil {
		fmt.Println(err)
		return
	}*/

	/*err := fillDatabaseWithSlaveServices()
	if err != nil {
		fmt.Println(err)
		return
	}*/
}
