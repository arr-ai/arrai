¼
PetShopModel«

PetShopModel"'
package"io.sysl.demo.petshop.model2
Employeeý
petshop.sysl:Ý
(
dob!`
petshop.sysl	
(
error
petshop.sysl
N

employeeId@B
patterns:
"pk
	"autoinc
petshop.sysl(
)
name!`
petshop.sysl

employeeId2®
Breed¤
petshop.sysl:
K
breedId@B
patterns:
"pk
	"autoinc
petshop.sysl		%
.
	breedName!`
petshop.sysl


,
species!`
petshop.sysl
,
numLegs!`
petshop.sysl
:
avgLifespan+R (`
petshop.sysl$
4
	avgWeight'R`
petshop.sysl 
,
legRank!`
petshop.sysl	
breedId2Ô
PetÌ
petshop.sysl:¬
(
dob!`
petshop.sysl	
,
numLegs!`
petshop.sysl
I
petId@B
patterns:
"pk
	"autoinc
petshop.sysl#
S
breedIdH
petshop.syslJ)


PetShopModelPetBreedbreedId
)
name!`
petshop.sysl
petId2½
EmployeeTendsPet¨
petshop.sysl:


employeeIdqB
patterns:
"pk
petshop.sysl.J<
"

PetShopModelEmployeeTendsPetEmployee
employeeId
p
petIdgB
patterns:
"pk
petshop.syslJ2
"

PetShopModelEmployeeTendsPetPetpetId

employeeId
petId
petshop.sysl2´
PetShopFacade¢

PetShopFacade"(
package"io.sysl.demo.petshop.facadeJH

PetShopModel2
Pet 2
EmployeeTendsPet 2
Employee 2	
Breed 
petshop.sysl4·


PetShopApi¨



PetShopApi"%
package"io.sysl.demo.petshop.api"
patterns
:
"rest*n
GET /petshop^
GET /petshop"
patterns
:
"rest:B	
PetShopB/petshop
petshop.sysl'*2Í
EmployeeÀ
petshop.sysl/4 
)
name!`
petshop.sysl00
(
dob!`
petshop.sysl11	
I
index@B
patterns:
"xml_attribute
petshop.sysl22$2æ
BreedÜ
petshop.sysl4;¼
I
index@B
patterns:
"xml_attribute
petshop.sysl99$
)
name!`
petshop.sysl55
,
species!`
petshop.sysl66
d
pets\
petshop.sysl77j=
petshop.sysl77J



PetShopApiBreedPet
0
avgLifespan!`
petshop.sysl882×
PetÏ
petshop.sysl;B¯
)
name!`
petshop.sysl<<
(
dob!`
petshop.sysl==	
,
numLegs!`
petshop.sysl>>
*
legRank
petshop.sysl??2¸
PetShop¬
petshop.sysl*/
p
	employeesc
petshop.sysl++jD
petshop.sysl++J%



PetShopApiPetShop
Employee
j
breeds`
petshop.sysl,,jA
petshop.sysl,,J"



PetShopApiPetShopBreed
,
numLegs!`
petshop.sysl--
petshop.sysl%%5Ì*
PetShopModelToApi¶*

PetShopModelToApi"'
package"io.sysl.demo.petshop.views2º
PetRankedByLeg§
petshop.syslpy
,
legRank!`
petshop.syslvv
(
petId
petshop.syslqq
*
breedId
petshop.syslrr
)
name!`
petshop.syslss
(
dob!`
petshop.sysltt	
,
numLegs!`
petshop.sysluuRÍ

modelToApi¾

petshopJ

PetShopModel2J0


PetShopModelToApi


PetShopApiPetShopj2J0


PetShopModelToApi


PetShopApiPetShop
petshop.syslDT	"É
&
petshop.syslDD
petshop.®«

rankedPetsj+j)J'%

PetShopModelToApiPetRankedByLegv
petshop.syslEEe"Zlet rankedPets = .table of Pet rank<PetShopModelToApi.PetRankedByLeg>(.numLegs as legRank)bóH
petshop.syslEE')
 
petshop.syslEE'
.Pet J
petshop.syslERES+
 
petshop.syslERES
.numLegsJ
petshop.syslERES+
 
petshop.syslERES
.numLegs".* 2legRank½
º
	employees¬I
petshop.syslGG8"-employees = employeeToApi(.table of Employee)2^
employeeToApiM
petshop.syslG&G0.
 
petshop.syslG&G0
.Employee ã
à
breedsÕL
petshop.syslHH;"0breeds = breedToApi(.table of Breed, rankedPets)2

breedToApiJ
petshop.syslH H*+
 
petshop.syslH H*
.Breed )
petshop.syslH1H1

rankedPets½
º
numLegs®F
petshop.syslII5"*numLegs = .table of Pet sum(.numLegs ?? 0)bâH
petshop.syslII )
 
petshop.syslII 
.Pet 
petshop.syslI(I4BqJ
petshop.syslI(I)+
 
petshop.syslI(I)
.numLegs!
petshop.syslI4I4 ".hf
pp`2
petshop.syslLL"let pp = .table of Pet)
 
petshop.syslLL
.Pet Á¾
bb·jr :
petshop.syslMM%"let bb = pp ~> .table of BreedBt!
petshop.syslMM
ppJ
petshop.syslMM%+
 
petshop.syslMM%
.Breed **Â¿
yy¸jr ;
petshop.syslNN&"let yy = pp !~> .table of BreedBt!
petshop.syslNN
ppJ
petshop.syslNN&+
 
petshop.syslNN&
.Breed **»¸
p²jr :
petshop.syslPP"let p = pp any(1) singleOrNull:ok
petshop.syslPP2L
.any!
petshop.syslPP
pp!
petshop.syslPPca
b\.
petshop.syslQQ"let b = p -> BreedJ)
 
petshop.syslQQ
pBreedom
pp2f6
petshop.syslRR#"let pp2 = b ?-> set of PetJ+
 
petshop.syslRR
bPet F
petshop.syslCV"*!view modelToApi(petshop <: PetShopModel):Rû
employeeToApié
,
employee jJ

PetShopModelEmployee5j3J1


PetShopModelToApi


PetShopApiEmployee¤j5j3J1


PetShopModelToApi


PetShopApiEmployee
petshop.syslW[	"Í
'
petshop.syslWW
employee.Q
O
nameG
petshop.syslXX(
 
petshop.syslXX
.nameO
M
dobF
petshop.syslYY'
 
petshop.syslYY
.dob{
y
indexp7
petshop.syslZZ&"index = autoinc("Employee")24
autoinc)
petshop.syslZZ
"EmployeeZ
petshop.syslV]">!view employeeToApi(employee <: set of PetShopModel.Employee):RÕ

breedToApiÆ
&
breedjJ

PetShopModelBreed
2
pet+j)J'%

PetShopModelToApiPetRankedByLeg2j0J.


PetShopModelToApi


PetShopApiBreed®j2j0J.


PetShopModelToApi


PetShopApiBreed
petshop.sysl^d	"Ú
$
petshop.sysl^^
breed.i
g
name_-
petshop.sysl__"name = .breedName-
 
petshop.sysl__
.	breedNameW
U
speciesJ
petshop.sysl``+
 
petshop.sysl``
.speciesý
ú
petsñH
petshop.syslaa7",pets = petToApi(-> set of Pet ~[petId]> pet)2£
petToApi
petshop.sysla*a4BwH
petshop.syslaa&J)
 
petshop.syslaa&
.Pet "
petshop.sysla4a4
pet*petIdt
r
avgLifespancj.
petshop.syslbb"avgLifespan = -1.0:,(j
petshop.syslbbZ1.0u
s
indexj4
petshop.syslcc#"index = autoinc("Breed")21
autoinc&
petshop.syslcc"Breed
petshop.sysl]f"e!view breedToApi(breed <: set of PetShopModel.Breed, pet <: set of PetShopModelToApi.PetRankedByLeg):Rñ
petToApiä
2
pet+j)J'%

PetShopModelToApiPetRankedByLeg0j.J,


PetShopModelToApi


PetShopApiPetj0j.J,


PetShopModelToApi


PetShopApiPet
petshop.syslgl	"Ë
"
petshop.syslgg
pet.Q
O
nameG
petshop.syslhh(
 
petshop.syslhh
.nameO
M
dobF
petshop.syslii'
 
petshop.syslii
.dobW
U
numLegsJ
petshop.sysljj+
 
petshop.sysljj
.numLegs¤
¡
legRank9
petshop.syslkk("legRank = fibonacci(.legRank)2W
	fibonacciJ
petshop.syslk k!+
 
petshop.syslk k!
.legRank[
petshop.syslfn"?!view petToApi(pet <: set of PetShopModelToApi.PetRankedByLeg):R
	fibonaccit

nB
patterns:

"abstractH
petshop.syslnn%",!view fibonacci(n <: int) -> int [~abstract]
petshop.syslBB7.
PetShopApiToModelò-

PetShopApiToModel"'
package"io.sysl.demo.petshop.views2©
AnonType_0__
I
breed@J>
#

PetShopApiToModelAnonType_0__

PetShopModelBreed
H
pets@j>J<
#

PetShopApiToModelAnonType_0__

PetShopModelPetRè+

apiToModelÙ+
&
petshopJ


PetShopApiPetShop+J)


PetShopApiToModel

PetShopModel³*j+J)


PetShopApiToModel

PetShopModel
petshop.sysl{«	"å)
&
petshop.sysl{{
petshop.
_breedsAndPetsj@j>J<


PetShopApiToModel#

PetShopApiToModelAnonType_0__
petshop.sysl|"çlet _breedsAndPets = .breeds -> <set of>(:
                let breedId = autoinc()

                breed = -> <PetShopModel.Breed>(:
                    breedId = breedId
                    species = .species
                    breedName = .name
                )

                pets = .pets -> <set of PetShopModel.Pet>(:
                    petId = autoinc()
                    breedId = breedId
                    .name
                    .dob
                )
            )
"´
I
petshop.sysl|!|"*
 
petshop.sysl|!|"
.breeds.NL
breedIdA3
petshop.sysl}}&"let breedId = autoinc()2	
autoincÝ
Ú
breedÐj2J0


PetShopApiToModel

PetShopModelBreedÅ
petshop.sysl"§breed = -> <PetShopModel.Breed>(:
                    breedId = breedId
                    species = .species
                    breedName = .name
                )
"Ð
!
petshop.sysl
..H
F
breedId;/
petshop.sysl"breedId = breedId
breedIdo
m
speciesb0
petshop.sysl"species = .species-
"
petshop.sysl
.speciesm
k
	breedName^/
petshop.sysl!"breedName = .name*
"
petshop.sysl !
.name³
°
pets§j2j0J.


PetShopApiToModel

PetShopModelPetÜ
petshop.sysl"½pets = .pets -> <set of PetShopModel.Pet>(:
                    petId = autoinc()
                    breedId = breedId
                    .name
                    .dob
                )
"
K
petshop.sysl*
"
petshop.sysl
.pets.H
F
petId=/
petshop.sysl$"petId = autoinc()2	
autoincH
F
breedId;/
petshop.sysl"breedId = breedId
breedIdU
S
nameK
petshop.sysl*
"
petshop.sysl
.nameS
Q
dobJ
petshop.sysl)
"
petshop.sysl
.dob
breedsAndPetsI
petshop.sysl/"+let breedsAndPets = _breedsAndPets snapshotb3/
petshop.sysl  
_breedsAndPets
minId1L
petshop.sysl9".let minId1 = breedsAndPets max(.breed.breedId)b¯.
petshop.sysl
breedsAndPetsx
petshop.sysl12W
L
petshop.sysl+,+
"
petshop.sysl+,
.breedbreedId".
minId2L
petshop.sysl9".let minId2 = breedsAndPets max(.breed.breedId)b¯.
petshop.sysl
breedsAndPetsx
petshop.sysl12W
L
petshop.sysl+,+
"
petshop.sysl+,
.breedbreedId".äá
leeÙj7j5J3


PetShopApiToModel

PetShopModelEmployee
petshop.sysl"alet lee = {{:}} -> <set of PetShopModel.Employee>(:
                name = "Bruce"
            )
"
Njj 
petshop.syslZ'
%j 
petshop.syslr .F
D
name<j,
petshop.sysl"name = "Bruce""Brucenl
leelessaj0
petshop.sysl"let leeless = !lee:($
petshop.sysl
leeLJ
decimal1>j0
petshop.sysl"let decimal1 = 1.2Z1.2LJ
decimal2>j0
petshop.sysl"let decimal2 = 2.3Z2.3ÜÙ
decimal3Ìjr A
petshop.sysl'"#let decimal3 = +decimal1 + decimal2BRj
petshop.sysl:-)
petshop.sysl
decimal1)
petshop.sysl''
decimal2«¨
	one_thirdj7
petshop.sysl""let one_third = 1.0 / 3.0BZ
*j
petshop.syslZ1.0*j
petshop.sysl""Z3.0­ª
	one_ninthj<
petshop.sysl)"let one_ninth = one_third ** 2BW*
petshop.sysl
	one_third'j
petshop.sysl))¦£
	bruce_leejr \
petshop.sysl¡¡I">let bruce_lee = lee singleOrNull?.name + substr(" Leek", 0, 4)B¯t
petshop.sysl¡,¡.S
I
petshop.sysl¡ ¡ :($
petshop.sysl¡¡
leename´3
petshop.sysl¡5¡I"substr(" Leek", 0, 4)2|
substr(
petshop.sysl¡<¡<" Leek#
petshop.sysl¡E¡E #
petshop.sysl¡H¡Hä
á
EmployeeÒj7j5J3


PetShopApiToModel

PetShopModelEmployee
petshop.sysl£ ¨"÷
P
petshop.sysl£ £!/
"
petshop.sysl£ £!
.	employees.R
P

employeeIdB4
petshop.sysl¤¤%"employeeId = autoinc()2	
autoincU
S
nameK
petshop.sysl¥¥*
"
petshop.sysl¥¥
.nameS
Q
dobJ
petshop.sysl¦¦)
"
petshop.sysl¦¦
.dob

errorjr 5
petshop.sysl§§!"error = minId2 - minId1BT'
petshop.sysl§§
minId2'
petshop.sysl§!§!
minId1h
f
Breed[
petshop.sysl©+©6J:
.
petshop.sysl©©
breedsAndPets.breed µ
²
Pet¨jr 
petshop.syslª)ª)B.
petshop.syslªª
breedsAndPetsK
petshop.syslª1ª2*
"
petshop.syslª1ª2
.pets".K
petshop.syslz¬"0!view apiToModel(petshop <: PetShopApi.PetShop):
petshop.syslyy7