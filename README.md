# website


Querys uitvoeren op poort 3306 (MySQL) voor de database voordat je de website test. P.S: alter tables werkt niet.


Je moet een config.json betsand hebben met je juiste credentials voor de database heb je die niet? vraag Bram Platier voor informatie



to-do
1. conformatie reserveren (mail sturen?)
2. google inloggen (rlly nice to have)



```
-- --------------------------------------------------------

--
-- Tabelstructuur voor tabel `accommodations`
--

CREATE TABLE `accommodations` (
  `ID` int(10) NOT NULL,
  `locatie` varchar(59) NOT NULL,
  `prijs` decimal(28,0) NOT NULL,
  `beschikbaar` tinyint(1) NOT NULL,
  `beschrijving` varchar(255) NOT NULL,
  `naam` varchar(90) NOT NULL,
  `imgurl` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Gegevens worden geëxporteerd voor tabel `accommodations`
--

INSERT INTO `accommodations` (`ID`, `locatie`, `prijs`, `beschikbaar`, `beschrijving`, `naam`, `imgurl`) VALUES
(2, 'Eindhoven', 10, 1, 'hihihihi', 'julian stinkt', 'images/vakantiehuisje.png'),
(3, 'Belgie', 100, 1, 'Een gezellig chalet in de natuur', 'Chalet', 'images/vakantiehuisje.png'),
(4, 'Duitsland', 200, 1, 'Luxe villa met zwembad', 'Villa', 'images/vakantiehuisje.png'),
(5, 'Nederland', 150, 1, 'Modern appartement in de stad', 'Appartement', 'images/vakantiehuisje.png'),
(6, 'Duitsland', 250, 1, 'Prachtig strandhuis met zeezicht', 'Strandhuis', 'images/vakantiehuisje.png'),
(7, 'Nederland', 180, 1, 'Knusse berghut in de bergen', 'Berghut', 'images/vakantiehuisje.png'),
(8, 'Belgie', 100, 1, 'Chalet met uitzicht op de bergen', 'Chalet', 'images/vakantiehuisje.png'),
(9, 'Duitsland', 200, 1, 'Moderne villa met tuin', 'Villa', 'images/vakantiehuisje.png'),
(10, 'Nederland', 150, 1, 'Ruim appartement in het centrum', 'Appartement', 'images/vakantiehuisje.png'),
(11, 'Duitsland', 250, 1, 'Strandhuis direct aan het strand', 'Strandhuis', 'images/vakantiehuisje.png'),
(12, 'Nederland', 180, 1, 'Berghut met open haard', 'Berghut', 'images/vakantiehuisje.png'),
(13, 'Belgie', 100, 1, 'Chalet omgeven door bossen', 'Chalet', 'images/vakantiehuisje.png'),
(14, 'Duitsland', 200, 1, 'Luxe villa met grote tuin', 'Villa', 'images/vakantiehuisje.png'),
(15, 'Nederland', 150, 1, 'Appartement met balkon', 'Appartement', 'images/vakantiehuisje.png'),
(16, 'Duitsland', 250, 1, 'Strandhuis met veranda', 'Strandhuis', 'images/vakantiehuisje.png'),
(17, 'Nederland', 180, 1, 'Berghut met prachtig uitzicht', 'Berghut', 'images/vakantiehuisje.png'),
(18, 'Belgie', 100, 1, 'Chalet aan de rivier', 'Chalet', 'images/vakantiehuisje.png'),
(19, 'Duitsland', 200, 1, 'Villa met binnenzwembad', 'Villa', 'images/vakantiehuisje.png'),
(20, 'Nederland', 150, 1, 'Appartement in rustige buurt', 'Appartement', 'images/vakantiehuisje.png'),
(21, 'Duitsland', 250, 1, 'Strandhuis met groot terras', 'Strandhuis', 'images/vakantiehuisje.png'),
(22, 'Nederland', 180, 1, 'Berghut met sauna', 'Berghut', 'images/vakantiehuisje.png');

-- --------------------------------------------------------

--
-- Tabelstructuur voor tabel `persoonsgegevens`
--

CREATE TABLE `persoonsgegevens` (
  `ID` int(11) NOT NULL,
  `voornaam` varchar(51) NOT NULL,
  `achternaam` varchar(58) NOT NULL,
  `email` varchar(254) NOT NULL,
  `kenteken` varchar(12) DEFAULT NULL,
  `wachtwoord` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Gegevens worden geëxporteerd voor tabel `persoonsgegevens`
--

INSERT INTO `persoonsgegevens` (`ID`, `voornaam`, `achternaam`, `email`, `kenteken`, `wachtwoord`) VALUES
(1, 'bram', 'platier', 'bramplatier@gmail.com', '63-20-bd', 'asdfghjkl'),
(2, 'Bram', 'Platier', 'PS220816@summacollege.nl', 'tew-22-2', '$2a$10$ItHAoPx7IflTUlAgbKJALe/au2iKDqdNxMwCZnG9lLEnc788SvLzm'),
(3, '', '', 'bramplatier@gmail.com', '', '$2a$10$MzhFcr.mnMEdir2ZUXWCo.zXXi0fcnzONpxoiV1OOFnlywsBNYGoi'),
(4, '', '', 'bramplatier@gmail.com', '', '$2a$10$QPY264svt4qfCyHmRAtIfuREA4Ge2D3Koex.Uo8flFS26wkCqSYMe'),
(5, 'Bram', 'Platier', 'bramp@gmail.com', 'tew-22-6', '$2a$10$XC2NLB0JxEuwFrjF.lvhY.MB5aMRsEw4V/Z2N2B0BaouvqUxZrnLC'),
(14, 'Bram', 'Platier', 'PS220816@summacollege.nl', 'sdf', '$2a$10$t4VvoigZP94SzJpOGMP6IOxF.h8AeJvVLzkmvk038qtu9YpbSiHFS'),
(17, 'marc', 'sdfs', 'PS220816@summacollege.nl', 'vnb', '$2a$10$JUm7i/GlMN506Za2yT/VOuq6Tng3OhWkC3kE1bPr2aOg0yoXS3rFq'),
(18, 'stijn', 'vdpol', 'test@gmail.com', 'vnb', '$2a$10$PoXIdYKXW6b58f5AvDV7aeIhXwRbz7uffY5F3/giznwOonKnUw1M2'),
(19, 't', 't', 't@t.com', 't', '$2a$10$AhP4JNYHi04OW2tmsX3cCe5CeMId5tGFd.b.ykEluWupPH97U7z5y'),
(21, 'Test', 'Test', 'test@test.nl', 'test001', '$2a$10$c0KcESR7M1lZxWQnJA4xyeYbxknIRICsBc58G1gSOqpoET7JqcZCi'),
(22, 'a', 'a', 'a@a.nl', 'a', '$2a$10$3YHkEfZh9gm0xMeFGhHRpeIrEOytAduhaCOiu/T/py.wUB1NoJAyq'),
(23, 'bramg', 'g', 'g@g.com', 'ggggg', '$2a$10$y9KT28cSJaPH43J84HK4W.ITZX6pjd3JIJjkhqBnn5wfoAgFPMthK');

-- --------------------------------------------------------

--
-- Tabelstructuur voor tabel `reservering`
--

CREATE TABLE `reservering` (
  `ID` int(10) NOT NULL,
  `locatie` varchar(30) NOT NULL,
  `email` varchar(254) NOT NULL,
  `indatum` date NOT NULL,
  `uitdatum` date NOT NULL,
  `naam` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Gegevens worden geëxporteerd voor tabel `reservering`
--

INSERT INTO `reservering` (`ID`, `locatie`, `email`, `indatum`, `uitdatum`, `naam`) VALUES
(22, 'Nederland', 'john@example.com', '2024-06-01', '2024-06-05', 'Villa Sunshine'),
(23, 'Duitsland', 'emma@example.com', '2024-06-02', '2024-06-08', 'Cottage Retreat'),
(27, 'België', 'ava@example.com', '2024-06-15', '2024-06-20', 'Lakeview Cottage'),
(28, 'Nederland', 'james@example.com', '2024-06-20', '2024-06-25', 'Riverside Retreat'),
(32, 'Duitsland', 'benjamin@example.com', '2024-07-01', '2024-07-05', 'Alpine Chalet'),
(33, 'België', 'mia@example.com', '2024-07-05', '2024-07-10', 'Desert Oasis'),
(34, 'Nederland', 'logan@example.com', '2024-07-08', '2024-07-15', 'Countryside Manor'),
(35, 'Duitsland', 'harper@example.com', '2024-07-10', '2024-07-18', 'Lakeside Cabin'),
(36, 'België', 'ethan@example.com', '2024-07-15', '2024-07-20', 'Riverfront House'),
(37, 'Nederland', 'ava@example.com', '2024-07-18', '2024-07-25', 'Tropical Bungalow'),
(38, 'Duitsland', 'olivia@example.com', '2024-07-20', '2024-07-28', 'Mountain Retreat'),
(39, 'België', 'mason@example.com', '2024-07-25', '2024-07-30', 'City Penthouse'),
(40, 'Nederland', 'william@example.com', '2024-07-28', '2024-08-03', 'Beach House'),
(41, 'Duitsland', 'emma@example.com', '2024-07-30', '2024-08-05', 'Lake House');

-- --------------------------------------------------------

--
-- Tabelstructuur voor tabel `slagboom`
--

CREATE TABLE `slagboom` (
  `kenteken` varchar(12) NOT NULL,
  `locatie` varchar(20) NOT NULL,
  `indatum` date NOT NULL,
  `uitdatum` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Gegevens worden geëxporteerd voor tabel `slagboom`
--

INSERT INTO `slagboom` (`kenteken`, `locatie`, `indatum`, `uitdatum`) VALUES
('63-20-bd', 'eindhoven', '2024-05-22', '2024-06-30');

--
-- Indexen voor geëxporteerde tabellen
--

--
-- Indexen voor tabel `accommodations`
--
ALTER TABLE `accommodations`
  ADD PRIMARY KEY (`ID`);

--
-- Indexen voor tabel `persoonsgegevens`
--
ALTER TABLE `persoonsgegevens`
  ADD PRIMARY KEY (`ID`);

--
-- Indexen voor tabel `reservering`
--
ALTER TABLE `reservering`
  ADD PRIMARY KEY (`ID`);

--
-- AUTO_INCREMENT voor geëxporteerde tabellen
--

--
-- AUTO_INCREMENT voor een tabel `accommodations`
--
ALTER TABLE `accommodations`
  MODIFY `ID` int(10) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;

--
-- AUTO_INCREMENT voor een tabel `persoonsgegevens`
--
ALTER TABLE `persoonsgegevens`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=24;

--
-- AUTO_INCREMENT voor een tabel `reservering`
--
ALTER TABLE `reservering`
  MODIFY `ID` int(10) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=42;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
```
