- Probabitlity Calulator : [Bloom Filter Probability Calculator](https://hur.st/bloomfilter/?n=1000000&p=0.01&m=10000000)


$$
p = \left(1 - \exp\left(-\frac{kn}{m}\right)\right)^{k}
$$

$$
m = \frac{-n \ln p}{(\ln 2)^2}
$$

$$
k = \frac{m}{n} \ln 2
$$

$$
n = \frac{m \ln p}{(\ln 2)^2}
$$