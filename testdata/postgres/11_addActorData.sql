-- Richer, IMDb-style bios (PostgreSQL-friendly using dollar-quoted strings)

UPDATE actor
SET bio = $$
Alfredo James “Al” Pacino was born on April 25, 1940, in East Harlem, New York City, to Sicilian Italian-American parents Rose Gerardi and Salvatore Pacino. After his parents split when he was young, he was raised largely in the South Bronx by his mother and grandparents, drifting between school and the streets while obsessing over movies and acting. He trained seriously at HB Studio and the Actors Studio, building a reputation in New York theatre before exploding on screen as Michael Corleone in The Godfather. Pacino became synonymous with volcanic intensity and precision—moving from Serpico and Dog Day Afternoon to Scarface, Heat, and later reinventions in films like Scent of a Woman (which earned him the Academy Award). Over decades he’s collected the industry’s top prizes, including Oscars, Tonys, and Emmys—cementing a career that bridges stage discipline with cinematic mythmaking.
$$
WHERE name = 'Al Pacino';

UPDATE actor
SET bio = $$
Robin McLaurin Williams was born July 21, 1951, in Chicago, Illinois, to Laurie McLaurin Williams, a former model from Mississippi, and Robert Fitzgerald Williams, a Ford executive. A shy kid with a fast mind, he found confidence through drama and improvisation, eventually studying in California and becoming a phenomenon in San Francisco’s comedy scene. He broke nationally as the extraterrestrial Mork on Mork & Mindy, then turned that manic speed into film stardom—Good Morning, Vietnam, Dead Poets Society, The Fisher King, Aladdin’s Genie, and Mrs. Doubtfire. Just as memorable was his dramatic work, culminating in an Academy Award for Good Will Hunting. Beloved for improvisational lightning and emotional warmth, he also battled addiction and health struggles in later years. Williams died on August 11, 2014, leaving a body of work that made laughter feel like something urgent and human.
$$
WHERE name = 'Robin Williams';

UPDATE actor
SET bio = $$
Kirsten Dunst was born April 30, 1982, in Point Pleasant, New Jersey, to Inez Rupprecht, a flight attendant and artist, and Klaus Dunst, a medical-services executive. She began working in commercials as a child and quickly graduated to films, displaying an uncommon poise for her age. Her breakout came when she stole scenes as the eerie, eternal Claudia in Interview with the Vampire, proving she could stand toe-to-toe with major stars. Dunst became a 1990s fixture with roles like Judy in Jumanji, then transitioned into adult leading parts—most famously as Mary Jane Watson in Spider-Man, while also earning acclaim in auteur-driven work that showcased her range and restraint. Known for balancing mainstream hits with riskier character pieces, she built a career around emotional clarity, dry intelligence, and an ability to make even heightened material feel lived-in.
$$
WHERE name = 'Kirsten Dunst';

UPDATE actor
SET bio = $$
Bonnie Hunt was born September 22, 1961, in Chicago, Illinois, one of several siblings in a big, working-class family. Before Hollywood, she lived a very un-Hollywood life: she studied and worked in health care, spending time as an oncology nurse while chasing comedy at night. That real-world grounding became part of her trademark—warmth, quick wit, and an improviser’s ease. She co-founded an improv troupe, performed with Chicago’s famed Second City, and landed a film break that turned into a run of memorable supporting roles: Rain Man, Beethoven, Jumanji, Jerry Maguire, The Green Mile, and more. Beyond acting, Hunt became a creator and host—writing, producing, and starring in television projects that leaned into her conversational, “it’s happening in the room” style. She’s also a sought-after voice performer, bringing that same lived-in humor to animated features.
$$
WHERE name = 'Bonnie Hunt';

UPDATE actor
SET bio = $$
Jack Lemmon was born February 8, 1925, in Newton, Massachusetts—famously, during the trip to the hospital—and grew up an only child with a restless curiosity and a knack for performance. He pursued acting early, sharpened his craft through theatre, and entered Hollywood as a seemingly effortless blend of charm and nerves. That “ordinary man in extraordinary trouble” persona became his superpower, whether he was sprinting through farce or breaking your heart in a drama. Two Academy Awards and a mountain of nominations traced his range: from Mister Roberts to The Apartment, Some Like It Hot, Days of Wine and Roses, and later late-career triumphs including his acclaimed partnership with Walter Matthau. Lemmon’s screen presence—busy, sincere, and deceptively complex—made him one of the defining American actors of the 20th century. He died on June 27, 2001.
$$
WHERE name = 'Jack Lemmon';

UPDATE actor
SET bio = $$
Walter Matthau was born October 1, 1920, on New York City’s Lower East Side, the son of Jewish immigrant parents who scraped by through hard times and instability. That early edge—streetwise, wary, and funny in self-defense—never left his work. He trained in theatre, served during World War II, and built a career as a character actor before becoming an unlikely leading man: rumpled, sharp, and magnetic. He won an Academy Award for The Fortune Cookie and became iconic through his prickly chemistry with Jack Lemmon, especially in The Odd Couple and their later comedies. Matthau could turn a single pause into a punchline, but he was also a formidable dramatic presence when the material demanded it. Off-screen he was known for candor and a gambler’s streak, a real-life flaw that mirrored the lovable messiness he often played. He died on July 1, 2000.
$$
WHERE name = 'Walter Matthau';

UPDATE actor
SET bio = $$
Ann-Margret Olsson was born April 28, 1941, in Valsjöbyn, Sweden, to Anna Regina and Carl Gustav Olsson. After the family emigrated to the United States, she grew up with dance and performance as a second language, training early and working relentlessly until she had the polish of a classic entertainer. She detonated into the 1960s pop-culture bloodstream with Bye Bye Birdie and became a quintessential screen presence—equal parts singer, dancer, comedian, and dramatic actress—most famously opposite Elvis Presley in Viva Las Vegas. Over the decades she refused to fossilize into nostalgia, taking on character roles and earning praise for her durability and craft. In Grumpier Old Men she returned as the glamorous, complicating force Ariel, using her star aura as storytelling texture rather than a spotlight. Ann-Margret’s career is a rare long arc: teen-idol electricity, serious acting credibility, and a performer’s discipline that never quit.
$$
WHERE name = 'Ann-Margret';

UPDATE actor
SET bio = $$
Whitney Elizabeth Houston was born August 9, 1963, in Newark, New Jersey, into a family where music was both profession and church. Her mother, Cissy Houston, was a celebrated gospel and soul singer; her father, John Russell Houston Jr., worked as an administrator and later managed aspects of the family’s careers. Whitney sang in church as a child, learned phrasing and power from gospel, then became a once-in-a-generation pop vocalist whose voice seemed engineered for radio history. Her film debut in The Bodyguard turned a superstar singer into a movie star overnight, and Waiting to Exhale expanded her screen identity beyond the ingénue glow. Behind the success, her personal life became tabloid fuel, and her later years were marked by public struggles. Houston died on February 11, 2012, but her cultural footprint—records, films, and a vocal standard few can touch—remains enormous.
$$
WHERE name = 'Whitney Houston';

UPDATE actor
SET bio = $$
Angela Evelyn Bassett was born August 16, 1958, in New York City, and raised largely in Florida by a mother who expected excellence and self-reliance. She carried that discipline into higher education, earning degrees at Yale and training as an actor with the seriousness of a scholar. Bassett’s breakthrough arrived as Tina Turner in What’s Love Got to Do with It, a performance built on physical transformation and emotional force—earning major awards attention and announcing a star with gravity. From there she became the rare actor who can anchor prestige drama and command blockbuster scale, moving between intimate character work and iconic presence. In modern pop culture she’s also celebrated for her regal authority in franchise roles, while continuing to collect acclaim for television performances that highlight her precision and intensity. Bassett’s signature is power with control: heat, intelligence, and a sense that every line is backed by a life.
$$
WHERE name = 'Angela Bassett';

UPDATE actor
SET bio = $$
Loretta Devine was born August 21, 1949, in Houston, Texas, and grew up in a large family where personality had to be big to be heard. She trained formally in speech and drama before carving her name into theatre history as an original Broadway cast member of Dreamgirls—work that gave her performance muscles you can’t fake. On screen she became an essential supporting force: capable of comedy, raw heartbreak, and righteous fury, often within the same scene. Waiting to Exhale turned her into a household face, and television expanded her reach through long runs and recurring roles, culminating in an Emmy-winning turn on Grey’s Anatomy. Devine’s gift is a deep, human immediacy—she plays warmth without softness and toughness without cliché. Whether she’s delivering a punchline or a confession, she has the rare ability to feel like someone you’ve actually met.
$$
WHERE name = 'Loretta Devine';

UPDATE actor
SET bio = $$
Stephen Glenn “Steve” Martin was born August 14, 1945, in Waco, Texas, and grew up in Southern California, where he mixed showmanship with obsession-level practice. He worked early jobs in entertainment (including magic and theme-park gigs), then turned stand-up into an art form—smart, absurd, and meticulously constructed. By the late 1970s he was a cultural meteor: sold-out arenas, comedy albums, and a persona that could be both the joke and the architect of the joke. Film stardom followed with The Jerk and a long stretch of hits that revealed surprising emotional range beneath the silliness. He’s also a writer and accomplished musician, especially as a banjo player in bluegrass, proving the “wild and crazy guy” was always a disciplined craftsman. In later years, he evolved into an elder statesman of comedy without losing the oddball edge that made him famous.
$$
WHERE name = 'Steve Martin';

UPDATE actor
SET bio = $$
Diane Keaton—born Diane Hall on January 5, 1946, in Los Angeles—grew up absorbing performance from everywhere, including the pageant-world theatricality that fascinated her mother. She studied acting, pushed her way into New York theatre, and adopted “Keaton” as a professional name. Early screen roles made her a defining face of 1970s American film: the offbeat intelligence of Annie Hall (which won her an Academy Award), the quiet steel of Kay in The Godfather films, and later dramatic showcases that proved she wasn’t just “quirky”—she was precise. Keaton’s style became its own cultural language: androgynous tailoring, ironic vulnerability, and comedy that never played dumb. She later expanded into directing, writing, photography, and an outspoken love of architecture and design. Keaton died on October 11, 2025, leaving behind a legacy that reshaped what leading women could look and sound like on screen.
$$
WHERE name = 'Diane Keaton';

UPDATE actor
SET bio = $$
Martin Hayter Short was born March 26, 1950, in Hamilton, Ontario, into a family where music, wit, and toughness coexisted. His mother was a concertmistress, and he grew up performing with the intensity of someone who understood that laughter can be both escape and craft. After university, he plunged into improv and sketch, rising through Second City and then detonating on SCTV and Saturday Night Live with a gallery of characters built from fearless commitment rather than punchline chasing. Short became a screen favorite in films like Three Amigos! and the Father of the Bride series, while also maintaining a serious theatre career that earned major awards. His comedic style is high-energy but not sloppy—every flail, squeal, and eyebrow is engineered. In recent years he’s enjoyed a major late-career surge, proving he’s not just a nostalgia act but an active master of the form.
$$
WHERE name = 'Martin Short';

UPDATE actor
SET bio = $$
Harrison Ford was born July 13, 1942, in Chicago, Illinois, to Dorothy Nidelman, a former radio actress, and John William “Christopher” Ford, an actor-turned-advertising executive. He wasn’t groomed for stardom—he drifted, found acting late, and spent years in unglamorous supporting roles before financial reality pushed him into carpentry. That detour became legend: the working actor literally building sets and doors while waiting for the right break. It arrived with American Graffiti and then the seismic one-two punch of Han Solo and Indiana Jones, roles that fused sarcasm, decency, and bruised heroism into modern myth. Ford kept expanding—Blade Runner, Witness (which brought major awards attention), and a long run as the adult, intelligent action lead. Outside film he’s known for aviation and conservation work, plus a lifelong preference for privacy. Even as his career accumulated honors, he remained the rare superstar who still feels like a guy you could meet at the hardware store.
$$
WHERE name = 'Harrison Ford';

UPDATE actor
SET bio = $$
Julia Karin Ormond was born January 4, 1965, in Epsom, Surrey, and found acting young, excelling in school productions before pursuing professional training. After early television work in the U.K., she arrived in 1990s cinema with a classic leading-lady aura—elegant, intelligent, and emotionally direct. Within a short stretch she anchored a run of prominent films, including Legends of the Fall, First Knight, and Sabrina, where her grounded sincerity kept romantic spectacle from floating away. Rather than staying in one lane, Ormond moved between period dramas, modern thrillers, and later television, often choosing roles with moral friction or psychological texture. Her career includes acclaimed work in prestige TV films and series, where she frequently plays women navigating power—sometimes wounded, sometimes dangerous, always composed. Ormond’s screen signature is quiet intensity: the sense that she’s thinking a few moves ahead, even when her characters say nothing at all.
$$
WHERE name = 'Julia Ormond';

UPDATE actor
SET bio = $$
Gregory Buck Kinnear was born June 17, 1963, in Logansport, Indiana, and spent parts of his youth around the world due to his family’s diplomatic connections. He entered entertainment through television, developing an affable, quick-on-his-feet persona that made him a natural host. The surprise was how well that ease translated to film: Kinnear brought charm without smugness and vulnerability without self-pity. His breakout dramatic credibility arrived with an Academy Award–nominated performance in As Good as It Gets, after which he became a go-to for smart supporting turns and romantic leads. He slipped comfortably between glossy studio projects (like Sabrina) and more character-driven work, often playing decent men with messy edges. Kinnear’s best performances lean into contradiction—confidence that cracks, humor that masks insecurity—making him a reliable anchor in ensembles where chemistry matters as much as plot.
$$
WHERE name = 'Greg Kinnear';

UPDATE actor
SET bio = $$
Jonathan Taylor Thomas was born Jonathan Taylor Weiss on September 8, 1981, in Bethlehem, Pennsylvania, and grew into fame so quickly that his teen years became a cultural event. He became a household name as Randy Taylor on Home Improvement, balancing sitcom timing with a believable, sarcastic intelligence that aged well with the audience. At the same time, he voiced young Simba in The Lion King, giving one of Disney’s most iconic heroes a playful, earnest core. Unlike many child stars, Thomas stepped back at the height of popularity, choosing education and privacy over constant visibility, later studying at universities including Columbia after time at Harvard. He returned occasionally for directing and guest appearances rather than a full-scale comeback. His career is a snapshot of 1990s stardom: talent, saturation, and then a deliberate refusal to be consumed by the machine that made him famous.
$$
WHERE name = 'Jonathan Taylor Thomas';

UPDATE actor
SET bio = $$
Brad Renfro was born July 25, 1982, in Knoxville, Tennessee, and was discovered as a pre-teen with no acting background—chosen specifically because he looked like a kid who’d already seen too much. That authenticity became his signature. He exploded onto the screen in The Client, holding his own with veteran actors and instantly becoming one of the era’s most rawly gifted young performers. Renfro followed with a string of films that leaned into danger and vulnerability—Tom and Huck, Sleepers, Apt Pupil, Bully, and Ghost World—often playing boys on the edge of adulthood, where charm and self-destruction are separated by a breath. Off-screen, legal troubles and addiction repeatedly interrupted his momentum, turning a promising trajectory into a struggle for survival. He died on January 15, 2008, at just 25. Renfro’s legacy is painful and undeniable: a natural actor whose honesty onscreen still feels like exposed nerve.
$$
WHERE name = 'Brad Renfro';

UPDATE actor
SET bio = $$
Eric Schweig was born in Inuvik, Northwest Territories, on June 19, 1967, and his early life was marked by dislocation: he was adopted as an infant and later spoke openly about the cultural separation and identity pressure that came with that experience. He eventually found an outlet in acting, where presence matters as much as polish. Schweig gained wide recognition as Uncas in The Last of the Mohicans, bringing physical strength and quiet romantic intensity to a role that required more than dialogue to land. He followed with film and television work that often drew on his ability to project steadiness, dignity, and simmering emotion—qualities that made him memorable even in supporting parts. Over time, Schweig became known not only for screen roles but also for the candor with which he discussed heritage, adoption, and the long echo of forced assimilation policies. His work carries that gravity: lived experience translated into performance.
$$
WHERE name = 'Eric Schweig';

UPDATE actor
SET bio = $$
Jean-Claude Van Damme was born Jean-Claude Camille François Van Varenberg on October 18, 1960, in Brussels, Belgium, to Eliana and Eugène Van Varenberg. A small child with big nervous energy, he was enrolled in martial arts young and attacked training with the hunger of someone building a new identity. He earned high ranks in karate and competed in kickboxing, while also studying ballet—an unexpected foundation that later showed up in his flexibility and control. Chasing movie stardom, he moved to the United States and hustled through tiny jobs until Bloodsport turned him into a global action brand: the perfect fusion of athletic spectacle and earnestness. Hits like Kickboxer and later studio vehicles cemented him as a uniquely physical star, capable of both brutality and vulnerability. His career has included reinventions and self-parody, but the core appeal remains: a performer whose body tells the story before his mouth ever does.
$$
WHERE name = 'Jean-Claude Van Damme';

UPDATE actor
SET bio = $$
Powers Allen Boothe was born June 1, 1948, in Snyder, Texas, raised around hard work and big skies, and carried that rugged authority into every role. He trained seriously, earning degrees in drama, and built his career through stage and television before becoming a go-to screen presence for men with gravity—dangerous, charismatic, or both. He won major recognition early with an Emmy-winning portrayal of cult leader Jim Jones in Guyana Tragedy, a performance that showcased his ability to play conviction as something terrifying. Boothe later became beloved for morally murky characters—most famously as saloon owner Cy Tolliver on Deadwood—where his gravelly voice and measured menace could turn simple dialogue into threat. He moved easily between film and TV, often elevating genre material with intelligence and bite. Boothe died on May 14, 2017, leaving behind a career of tough-guy roles that were never just tough.
$$
WHERE name = 'Powers Boothe';

UPDATE actor
SET bio = $$
Raymond J. Barry was born March 14, 1939, in New York and developed as a performer with a stage actor’s seriousness, bringing that discipline into decades of film and television work. He became the kind of actor directors love: instantly believable, emotionally grounded, and capable of suggesting an entire backstory in a glance. Barry is best known for authoritative, high-pressure characters—politicians, military men, hard-edged professionals—often playing power as something brittle, insecure, or quietly cruel. His filmography spans studio pictures and independent work, and he’s a familiar face across modern television, where his calm delivery can make even a small role feel consequential. What sets him apart is restraint: he rarely “pushes” a scene, but the tension rises anyway because he understands how real people threaten one another—with tone, timing, and certainty rather than volume. In ensemble casts, he’s often the silent engine of unease.
$$
WHERE name = 'Raymond J. Barry';
