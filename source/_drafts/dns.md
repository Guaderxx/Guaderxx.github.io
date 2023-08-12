---
title: dns
tags: dns
categories: Protocol
copyright: Guader
date: 2023-06-23 14:31:02
keywords:
copyright_author_href:
copyright_info:
---


## Summary

The Domain Name System (DNS) is a hierarchical and distributed naming system for computers, services, and other resources in the Internet or other Internet Protocol (IP) networks.

The Internet maintains two principal namespaces
- the domain name hierarchy
- the IP address spaces

Internet name servers and a communication protocol implement the Domain Name System. A DNS name server is a server that stores the DNS records for a domain; a DNS name server responds with answers to queries against its database.

The most common types of records stored in the DNS database are:
- Start of authority                  SOA
- IP addresses                        A and AAAA
- SMTP mail exchangers                MX
- Name servers                        NS
- Pointers for reverse DNS lookups    PTR
- Domain name aliases                 CNAME

As a general purpose database, the DNS has also been used in combating unsolicited email (spam) by storing a real-time blockhole list (RBL). The DNS database is traditionally stored in a structed text file, the zone file, but other database systems are common.

The Domain Name System originally used the User Datagram Protocol(UDP) as transport over IP. Reliability, security, and privacy concerns spawned the use of the Transmission Control Protocol(TCP) as well as numerous other protocol developments.



## Function

An ofter-used analogy to explain the DNS is that it serves as the phone book for the Internet by translating human-friendly computer hostnames into IP addresses.

For example, the hostname `www.example.com` within the domain name `example.com` translates to the addresses 93.184.216.34 (IPv4) and 2606:2800:220:1:248:1893:25c8:1946 (IPv6).

An impoprtant and ubiquitous function of the DNS is its central role in distributed Internet services such as cloud services and content delivery networkds. When a user accesses a distributed Internet service using a URL, the domain name of the URL is translated to the IP address of a server that is proximal to the user. The key functionality of the DNS exploited here is that different users can simultaneously receive different translations for the same domain name, a key point of divergence from a traditional phone-book view of the DNS. This process of using the DNS to assign proximal servers to users is key to providing faster and more reliable responses on the Internet and is widely used by most major Internet services.



## Structure

### Domain name space

The domain name space consists of a tree data structure. Each node or leaf in the tree has a label and zero or more resource records(RR), which hold information associated with the domain name. The domain name itself consists of the label, concatenated with the name of its parent node on the right, separated by a dot.

The tree sub-divides into zones beginning at the root zone. A DNS zone may consist of as many domains and sub domains as the zone manager chooses. DNS can also be partitioned according to class where the separate classes can be thought of as an array of parallel namespace trees.


### Domain name syntax, internationalization

The definitive descriptions of the rules for forming domain names appear in RFC 1035, RFC 1123, RFC 2181, and RFC 5892. A domain name consists of one or more parts, technically called labels, that are conventionally concatenated, and delimited by dots, such as example.com.

The right-most label conveys the top-level domain; for example, the domain name www.example.com belongs to the top-level domain com.

The hierarchy of domains descends from right to left; each label to the left specifies a subdivision, or subdomain of the domain to the right. For example, the label example specifies a subdomain of the com domain, and `www` is a subdomain of example.com. The tree of subdivisions may have up to 127 levels.

A label may contain zeor to 63 characters. The null label, of length zero, is reserved for the root zone. The full domain name may not exceed the length of 253 characters in its textual representaion. In the internal binary representation of the DNS the maximum length requires 255 octets of storage, as it also stores the length of the name.


### Name servers

The Domain Name System is maintained by a distributed database system, which uses the client-server model. The nodes of this database are the name servers. Each domain has at least one authoritative DNS server that publishes information about that domain and the name servers of any domain subordinate to it. The top of the hierarchy is served by the root name servers, the servers to query when looking up (resolving) a TLD.


#### Authoritative name server

An authoritative name server is a name server that only gives answers to DNS queries from data that have been configured by an original source, for example, the domain administrator or by dynamic DNS methods, in contract to answers obtained via a query to another name server that only maintains a cache of data.


## Operation

### Address resolution mechanism



### DNS resolvers

The client side of the DNS is called a DNS resolver. A resolver is responsible for initiating and sequencing the queries that ultimately lead to a full resolution(translation) of the resource sought, e.g., translation of a domain name into an IP address. DNS resolvers are classified by a variety of query methods, such as recursive, non-recursive, and iterative. A resolution process may use a combination if these methods.

In a non-recursive query, a DNS resolver queries a DNS server that provides a record either for which the server is authoritative, or it provides a partial result without querying other servers. In case of a caching DNS resolver, the non-recursive query of its local DNS cache delivers a result and resuces the load on upstream DNS servers by caching DNS resource records for a period of time after an initial response from upstream DNS servers.

In a recursive query, a DNS resolver queries a single DNS server, which may in turn query other DNS servers on behalf of the requester. For example, a simple stub resolver running on a home router typically makes a recursive query to the DNS server run by the user's ISP. A recursive query is one for which the DNS server answers the query completely by querying other name servers as needed. In typical operation, a client issues a recursive query to a caching recursvie DNS server, which subsequently issues non-recursive queries to determine the answer and send a single answer back to the client.
