PGDMP                         x            db_hr    12.1 (Debian 12.1-1.pgdg90+1)    12.1 (Debian 12.1-1.pgdg90+1) �    w           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            x           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            y           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            z           1262    16384    db_hr    DATABASE     w   CREATE DATABASE db_hr WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'zh_CN.UTF-8' LC_CTYPE = 'zh_CN.UTF-8';
    DROP DATABASE db_hr;
                postgres    false            �            1259    16440 
   group_tags    TABLE     ]   CREATE TABLE public.group_tags (
    tag_id bigint NOT NULL,
    group_id bigint NOT NULL
);
    DROP TABLE public.group_tags;
       public         heap    postgres    false            �            1259    16390    permissions_roles    TABLE     k   CREATE TABLE public.permissions_roles (
    role_id bigint NOT NULL,
    permissions_id bigint NOT NULL
);
 %   DROP TABLE public.permissions_roles;
       public         heap    postgres    false            �            1259    16492    profile_groups    TABLE     e   CREATE TABLE public.profile_groups (
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL
);
 "   DROP TABLE public.profile_groups;
       public         heap    postgres    false            �            1259    16435    profile_tags    TABLE     a   CREATE TABLE public.profile_tags (
    tag_id bigint NOT NULL,
    profile_id bigint NOT NULL
);
     DROP TABLE public.profile_tags;
       public         heap    postgres    false            �            1259    16497 
   tag_groups    TABLE     ]   CREATE TABLE public.tag_groups (
    group_id bigint NOT NULL,
    tag_id bigint NOT NULL
);
    DROP TABLE public.tag_groups;
       public         heap    postgres    false            �            1259    16579    tb_audit    TABLE     �  CREATE TABLE public.tb_audit (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    operator_id bigint,
    auditor_id bigint,
    object text,
    action integer,
    fields text,
    org_object_id integer[],
    dest_object_id integer[],
    state integer,
    reply text,
    body text,
    remark text
);
    DROP TABLE public.tb_audit;
       public         heap    postgres    false            �            1259    16577    tb_audit_id_seq    SEQUENCE     x   CREATE SEQUENCE public.tb_audit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.tb_audit_id_seq;
       public          postgres    false    236            {           0    0    tb_audit_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.tb_audit_id_seq OWNED BY public.tb_audit.id;
          public          postgres    false    235            �            1259    16540    tb_group_transfer    TABLE     #  CREATE TABLE public.tb_group_transfer (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    profile bigint NOT NULL,
    old_group bigint NOT NULL,
    new_group bigint NOT NULL,
    new_group_combination integer[],
    description text,
    added_tags_record integer[]
);
 %   DROP TABLE public.tb_group_transfer;
       public         heap    postgres    false            �            1259    16538    tb_group_transfer_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tb_group_transfer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 /   DROP SEQUENCE public.tb_group_transfer_id_seq;
       public          postgres    false    230            |           0    0    tb_group_transfer_id_seq    SEQUENCE OWNED BY     U   ALTER SEQUENCE public.tb_group_transfer_id_seq OWNED BY public.tb_group_transfer.id;
          public          postgres    false    229            �            1259    16509 	   tb_groups    TABLE     t  CREATE TABLE public.tb_groups (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    name text NOT NULL,
    code integer DEFAULT 0,
    parent bigint,
    levels text,
    coefficient numeric DEFAULT 0,
    locked boolean,
    invalid boolean,
    is_default boolean
);
    DROP TABLE public.tb_groups;
       public         heap    postgres    false            �            1259    16507    tb_groups_id_seq    SEQUENCE     y   CREATE SEQUENCE public.tb_groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 '   DROP SEQUENCE public.tb_groups_id_seq;
       public          postgres    false    225            }           0    0    tb_groups_id_seq    SEQUENCE OWNED BY     E   ALTER SEQUENCE public.tb_groups_id_seq OWNED BY public.tb_groups.id;
          public          postgres    false    224            �            1259    24600 
   tb_message    TABLE     &  CREATE TABLE public.tb_message (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    text_id bigint NOT NULL,
    rec_id bigint NOT NULL,
    status integer NOT NULL,
    message_type text
);
    DROP TABLE public.tb_message;
       public         heap    postgres    false            �            1259    24598    tb_message_id_seq    SEQUENCE     z   CREATE SEQUENCE public.tb_message_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 (   DROP SEQUENCE public.tb_message_id_seq;
       public          postgres    false    248            ~           0    0    tb_message_id_seq    SEQUENCE OWNED BY     G   ALTER SEQUENCE public.tb_message_id_seq OWNED BY public.tb_message.id;
          public          postgres    false    247            �            1259    24609    tb_message_text    TABLE     �   CREATE TABLE public.tb_message_text (
    id bigint NOT NULL,
    send_id bigint NOT NULL,
    title text NOT NULL,
    text text NOT NULL,
    message_type text NOT NULL,
    "group" bigint,
    "postDate" timestamp with time zone,
    role bigint
);
 #   DROP TABLE public.tb_message_text;
       public         heap    postgres    false            �            1259    24607    tb_message_text_id_seq    SEQUENCE        CREATE SEQUENCE public.tb_message_text_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 -   DROP SEQUENCE public.tb_message_text_id_seq;
       public          postgres    false    250                       0    0    tb_message_text_id_seq    SEQUENCE OWNED BY     Q   ALTER SEQUENCE public.tb_message_text_id_seq OWNED BY public.tb_message_text.id;
          public          postgres    false    249            �            1259    16460    tb_operate_record    TABLE     {   CREATE TABLE public.tb_operate_record (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    body text
);
 %   DROP TABLE public.tb_operate_record;
       public         heap    postgres    false            �            1259    16458    tb_operate_record_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tb_operate_record_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 /   DROP SEQUENCE public.tb_operate_record_id_seq;
       public          postgres    false    216            �           0    0    tb_operate_record_id_seq    SEQUENCE OWNED BY     U   ALTER SEQUENCE public.tb_operate_record_id_seq OWNED BY public.tb_operate_record.id;
          public          postgres    false    215            �            1259    16603    tb_permissions    TABLE     N   CREATE TABLE public.tb_permissions (
    id bigint NOT NULL,
    name text
);
 "   DROP TABLE public.tb_permissions;
       public         heap    postgres    false            �            1259    16601    tb_permissions_id_seq    SEQUENCE     ~   CREATE SEQUENCE public.tb_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 ,   DROP SEQUENCE public.tb_permissions_id_seq;
       public          postgres    false    240            �           0    0    tb_permissions_id_seq    SEQUENCE OWNED BY     O   ALTER SEQUENCE public.tb_permissions_id_seq OWNED BY public.tb_permissions.id;
          public          postgres    false    239            �            1259    16551 
   tb_profile    TABLE     �  CREATE TABLE public.tb_profile (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    name text NOT NULL,
    job_number text,
    type_card text,
    phone text,
    id_card text NOT NULL,
    gender text,
    birth_day text,
    source text,
    school text,
    graduation_date text,
    specialty text,
    last_company text,
    first_job_date text,
    work_age integer DEFAULT 0,
    nation text,
    marital_status text,
    account_location text,
    address text,
    bank_card text,
    on_board_date text,
    freezed boolean,
    audit_state integer
);
    DROP TABLE public.tb_profile;
       public         heap    postgres    false            �            1259    16549    tb_profile_id_seq    SEQUENCE     z   CREATE SEQUENCE public.tb_profile_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 (   DROP SEQUENCE public.tb_profile_id_seq;
       public          postgres    false    232            �           0    0    tb_profile_id_seq    SEQUENCE OWNED BY     G   ALTER SEQUENCE public.tb_profile_id_seq OWNED BY public.tb_profile.id;
          public          postgres    false    231            �            1259    16409 	   tb_record    TABLE     �   CREATE TABLE public.tb_record (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    body text,
    object text
);
    DROP TABLE public.tb_record;
       public         heap    postgres    false            �            1259    16407    tb_record_id_seq    SEQUENCE     y   CREATE SEQUENCE public.tb_record_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 '   DROP SEQUENCE public.tb_record_id_seq;
       public          postgres    false    207            �           0    0    tb_record_id_seq    SEQUENCE OWNED BY     E   ALTER SEQUENCE public.tb_record_id_seq OWNED BY public.tb_record.id;
          public          postgres    false    206            �            1259    16397    tb_roles    TABLE     �   CREATE TABLE public.tb_roles (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    name text
);
    DROP TABLE public.tb_roles;
       public         heap    postgres    false            �            1259    16395    tb_roles_id_seq    SEQUENCE     x   CREATE SEQUENCE public.tb_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.tb_roles_id_seq;
       public          postgres    false    205            �           0    0    tb_roles_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.tb_roles_id_seq OWNED BY public.tb_roles.id;
          public          postgres    false    204            �            1259    16625 	   tb_salary    TABLE     �   CREATE TABLE public.tb_salary (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    template_account bigint NOT NULL,
    template text NOT NULL,
    year text NOT NULL,
    month text NOT NULL,
    locked boolean,
    data jsonb
);
    DROP TABLE public.tb_salary;
       public         heap    postgres    false            �            1259    16636    tb_salary_config    TABLE        CREATE TABLE public.tb_salary_config (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    base numeric,
    tax_threshold numeric,
    reference text
);
 $   DROP TABLE public.tb_salary_config;
       public         heap    postgres    false            �            1259    16634    tb_salary_config_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tb_salary_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 .   DROP SEQUENCE public.tb_salary_config_id_seq;
       public          postgres    false    246            �           0    0    tb_salary_config_id_seq    SEQUENCE OWNED BY     S   ALTER SEQUENCE public.tb_salary_config_id_seq OWNED BY public.tb_salary_config.id;
          public          postgres    false    245            �            1259    16614    tb_salary_fields    TABLE     �  CREATE TABLE public.tb_salary_fields (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    profile_id bigint,
    salary_id bigint,
    department_group_id bigint,
    post_group_id bigint,
    key text,
    name text,
    alias text,
    value numeric,
    content text,
    should_tax boolean,
    is_income boolean,
    is_deduct boolean,
    year text,
    month text,
    fit_into_year text,
    fit_into_month text
);
 $   DROP TABLE public.tb_salary_fields;
       public         heap    postgres    false            �            1259    16612    tb_salary_fields_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tb_salary_fields_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 .   DROP SEQUENCE public.tb_salary_fields_id_seq;
       public          postgres    false    242            �           0    0    tb_salary_fields_id_seq    SEQUENCE OWNED BY     S   ALTER SEQUENCE public.tb_salary_fields_id_seq OWNED BY public.tb_salary_fields.id;
          public          postgres    false    241            �            1259    16623    tb_salary_id_seq    SEQUENCE     y   CREATE SEQUENCE public.tb_salary_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 '   DROP SEQUENCE public.tb_salary_id_seq;
       public          postgres    false    244            �           0    0    tb_salary_id_seq    SEQUENCE OWNED BY     E   ALTER SEQUENCE public.tb_salary_id_seq OWNED BY public.tb_salary.id;
          public          postgres    false    243            �            1259    16420    tb_salary_profile_config    TABLE     5  CREATE TABLE public.tb_salary_profile_config (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    template_field_id text,
    profile_id bigint,
    operate text,
    value numeric,
    description text
);
 ,   DROP TABLE public.tb_salary_profile_config;
       public         heap    postgres    false            �            1259    16418    tb_salary_profile_config_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tb_salary_profile_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 6   DROP SEQUENCE public.tb_salary_profile_config_id_seq;
       public          postgres    false    209            �           0    0    tb_salary_profile_config_id_seq    SEQUENCE OWNED BY     c   ALTER SEQUENCE public.tb_salary_profile_config_id_seq OWNED BY public.tb_salary_profile_config.id;
          public          postgres    false    208            �            1259    16447    tb_tags    TABLE     �   CREATE TABLE public.tb_tags (
    id bigint NOT NULL,
    name text NOT NULL,
    coefficient numeric DEFAULT 0,
    parent bigint DEFAULT 0,
    commensalism_group_ids integer[]
);
    DROP TABLE public.tb_tags;
       public         heap    postgres    false            �            1259    16445    tb_tags_id_seq    SEQUENCE     w   CREATE SEQUENCE public.tb_tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 %   DROP SEQUENCE public.tb_tags_id_seq;
       public          postgres    false    214            �           0    0    tb_tags_id_seq    SEQUENCE OWNED BY     A   ALTER SEQUENCE public.tb_tags_id_seq OWNED BY public.tb_tags.id;
          public          postgres    false    213            �            1259    16528    tb_template    TABLE     i  CREATE TABLE public.tb_template (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    name text NOT NULL,
    type text,
    months integer,
    startup boolean,
    init_data text,
    "order" integer,
    user_id bigint NOT NULL,
    audit_state integer
);
    DROP TABLE public.tb_template;
       public         heap    postgres    false            �            1259    16566    tb_template_account    TABLE     �   CREATE TABLE public.tb_template_account (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    name text NOT NULL,
    "order" integer[]
);
 '   DROP TABLE public.tb_template_account;
       public         heap    postgres    false            �            1259    16564    tb_template_account_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tb_template_account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 1   DROP SEQUENCE public.tb_template_account_id_seq;
       public          postgres    false    234            �           0    0    tb_template_account_id_seq    SEQUENCE OWNED BY     Y   ALTER SEQUENCE public.tb_template_account_id_seq OWNED BY public.tb_template_account.id;
          public          postgres    false    233            �            1259    16526    tb_template_id_seq    SEQUENCE     {   CREATE SEQUENCE public.tb_template_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 )   DROP SEQUENCE public.tb_template_id_seq;
       public          postgres    false    228            �           0    0    tb_template_id_seq    SEQUENCE OWNED BY     I   ALTER SEQUENCE public.tb_template_id_seq OWNED BY public.tb_template.id;
          public          postgres    false    227            �            1259    16591    tb_usergroups    TABLE     �   CREATE TABLE public.tb_usergroups (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    name text NOT NULL,
    parent bigint,
    levels text
);
 !   DROP TABLE public.tb_usergroups;
       public         heap    postgres    false            �            1259    16589    tb_usergroups_id_seq    SEQUENCE     }   CREATE SEQUENCE public.tb_usergroups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 +   DROP SEQUENCE public.tb_usergroups_id_seq;
       public          postgres    false    238            �           0    0    tb_usergroups_id_seq    SEQUENCE OWNED BY     M   ALTER SEQUENCE public.tb_usergroups_id_seq OWNED BY public.tb_usergroups.id;
          public          postgres    false    237            �            1259    16476    tb_users    TABLE     G  CREATE TABLE public.tb_users (
    id bigint NOT NULL,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    email text,
    username text NOT NULL,
    nichname text NOT NULL,
    id_card text NOT NULL,
    password text NOT NULL,
    is_super boolean,
    picture text,
    state integer
);
    DROP TABLE public.tb_users;
       public         heap    postgres    false            �            1259    16474    tb_users_id_seq    SEQUENCE     x   CREATE SEQUENCE public.tb_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.tb_users_id_seq;
       public          postgres    false    219            �           0    0    tb_users_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.tb_users_id_seq OWNED BY public.tb_users.id;
          public          postgres    false    218            �            1259    16502    templateaccount_groups    TABLE     v   CREATE TABLE public.templateaccount_groups (
    group_id bigint NOT NULL,
    template_account_id bigint NOT NULL
);
 *   DROP TABLE public.templateaccount_groups;
       public         heap    postgres    false            �            1259    16521    templateaccount_templates    TABLE     |   CREATE TABLE public.templateaccount_templates (
    template_id bigint NOT NULL,
    template_account_id bigint NOT NULL
);
 -   DROP TABLE public.templateaccount_templates;
       public         heap    postgres    false            �            1259    16487    user_groups    TABLE     _   CREATE TABLE public.user_groups (
    group_id bigint NOT NULL,
    user_id bigint NOT NULL
);
    DROP TABLE public.user_groups;
       public         heap    postgres    false            �            1259    16385 
   user_roles    TABLE     ]   CREATE TABLE public.user_roles (
    role_id bigint NOT NULL,
    user_id bigint NOT NULL
);
    DROP TABLE public.user_roles;
       public         heap    postgres    false            �            1259    16430 	   user_tags    TABLE     [   CREATE TABLE public.user_tags (
    tag_id bigint NOT NULL,
    user_id bigint NOT NULL
);
    DROP TABLE public.user_tags;
       public         heap    postgres    false            �            1259    16469    user_usergroups    TABLE     h   CREATE TABLE public.user_usergroups (
    user_id bigint NOT NULL,
    user_group_id bigint NOT NULL
);
 #   DROP TABLE public.user_usergroups;
       public         heap    postgres    false            s           2604    16582    tb_audit id    DEFAULT     j   ALTER TABLE ONLY public.tb_audit ALTER COLUMN id SET DEFAULT nextval('public.tb_audit_id_seq'::regclass);
 :   ALTER TABLE public.tb_audit ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    235    236    236            o           2604    16543    tb_group_transfer id    DEFAULT     |   ALTER TABLE ONLY public.tb_group_transfer ALTER COLUMN id SET DEFAULT nextval('public.tb_group_transfer_id_seq'::regclass);
 C   ALTER TABLE public.tb_group_transfer ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    229    230    230            k           2604    16512    tb_groups id    DEFAULT     l   ALTER TABLE ONLY public.tb_groups ALTER COLUMN id SET DEFAULT nextval('public.tb_groups_id_seq'::regclass);
 ;   ALTER TABLE public.tb_groups ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    224    225    225            y           2604    24603    tb_message id    DEFAULT     n   ALTER TABLE ONLY public.tb_message ALTER COLUMN id SET DEFAULT nextval('public.tb_message_id_seq'::regclass);
 <   ALTER TABLE public.tb_message ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    248    247    248            z           2604    24612    tb_message_text id    DEFAULT     x   ALTER TABLE ONLY public.tb_message_text ALTER COLUMN id SET DEFAULT nextval('public.tb_message_text_id_seq'::regclass);
 A   ALTER TABLE public.tb_message_text ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    250    249    250            i           2604    16463    tb_operate_record id    DEFAULT     |   ALTER TABLE ONLY public.tb_operate_record ALTER COLUMN id SET DEFAULT nextval('public.tb_operate_record_id_seq'::regclass);
 C   ALTER TABLE public.tb_operate_record ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    216    215    216            u           2604    16606    tb_permissions id    DEFAULT     v   ALTER TABLE ONLY public.tb_permissions ALTER COLUMN id SET DEFAULT nextval('public.tb_permissions_id_seq'::regclass);
 @   ALTER TABLE public.tb_permissions ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    239    240    240            p           2604    16554    tb_profile id    DEFAULT     n   ALTER TABLE ONLY public.tb_profile ALTER COLUMN id SET DEFAULT nextval('public.tb_profile_id_seq'::regclass);
 <   ALTER TABLE public.tb_profile ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    231    232    232            d           2604    16412    tb_record id    DEFAULT     l   ALTER TABLE ONLY public.tb_record ALTER COLUMN id SET DEFAULT nextval('public.tb_record_id_seq'::regclass);
 ;   ALTER TABLE public.tb_record ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    207    206    207            c           2604    16400    tb_roles id    DEFAULT     j   ALTER TABLE ONLY public.tb_roles ALTER COLUMN id SET DEFAULT nextval('public.tb_roles_id_seq'::regclass);
 :   ALTER TABLE public.tb_roles ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    204    205    205            w           2604    16628    tb_salary id    DEFAULT     l   ALTER TABLE ONLY public.tb_salary ALTER COLUMN id SET DEFAULT nextval('public.tb_salary_id_seq'::regclass);
 ;   ALTER TABLE public.tb_salary ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    243    244    244            x           2604    16639    tb_salary_config id    DEFAULT     z   ALTER TABLE ONLY public.tb_salary_config ALTER COLUMN id SET DEFAULT nextval('public.tb_salary_config_id_seq'::regclass);
 B   ALTER TABLE public.tb_salary_config ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    245    246    246            v           2604    16617    tb_salary_fields id    DEFAULT     z   ALTER TABLE ONLY public.tb_salary_fields ALTER COLUMN id SET DEFAULT nextval('public.tb_salary_fields_id_seq'::regclass);
 B   ALTER TABLE public.tb_salary_fields ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    241    242    242            e           2604    16423    tb_salary_profile_config id    DEFAULT     �   ALTER TABLE ONLY public.tb_salary_profile_config ALTER COLUMN id SET DEFAULT nextval('public.tb_salary_profile_config_id_seq'::regclass);
 J   ALTER TABLE public.tb_salary_profile_config ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    209    208    209            f           2604    16450 
   tb_tags id    DEFAULT     h   ALTER TABLE ONLY public.tb_tags ALTER COLUMN id SET DEFAULT nextval('public.tb_tags_id_seq'::regclass);
 9   ALTER TABLE public.tb_tags ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    214    213    214            n           2604    16531    tb_template id    DEFAULT     p   ALTER TABLE ONLY public.tb_template ALTER COLUMN id SET DEFAULT nextval('public.tb_template_id_seq'::regclass);
 =   ALTER TABLE public.tb_template ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    228    227    228            r           2604    16569    tb_template_account id    DEFAULT     �   ALTER TABLE ONLY public.tb_template_account ALTER COLUMN id SET DEFAULT nextval('public.tb_template_account_id_seq'::regclass);
 E   ALTER TABLE public.tb_template_account ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    234    233    234            t           2604    16594    tb_usergroups id    DEFAULT     t   ALTER TABLE ONLY public.tb_usergroups ALTER COLUMN id SET DEFAULT nextval('public.tb_usergroups_id_seq'::regclass);
 ?   ALTER TABLE public.tb_usergroups ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    238    237    238            j           2604    16479    tb_users id    DEFAULT     j   ALTER TABLE ONLY public.tb_users ALTER COLUMN id SET DEFAULT nextval('public.tb_users_id_seq'::regclass);
 :   ALTER TABLE public.tb_users ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    219    218    219            N          0    16440 
   group_tags 
   TABLE DATA           6   COPY public.group_tags (tag_id, group_id) FROM stdin;
    public          postgres    false    212   E�       E          0    16390    permissions_roles 
   TABLE DATA           D   COPY public.permissions_roles (role_id, permissions_id) FROM stdin;
    public          postgres    false    203   b�       W          0    16492    profile_groups 
   TABLE DATA           >   COPY public.profile_groups (group_id, profile_id) FROM stdin;
    public          postgres    false    221   �       M          0    16435    profile_tags 
   TABLE DATA           :   COPY public.profile_tags (tag_id, profile_id) FROM stdin;
    public          postgres    false    211   ��       X          0    16497 
   tag_groups 
   TABLE DATA           6   COPY public.tag_groups (group_id, tag_id) FROM stdin;
    public          postgres    false    222   -�       f          0    16579    tb_audit 
   TABLE DATA           �   COPY public.tb_audit (id, "createdAt", "updatedAt", "deletedAt", operator_id, auditor_id, object, action, fields, org_object_id, dest_object_id, state, reply, body, remark) FROM stdin;
    public          postgres    false    236   O�       `          0    16540    tb_group_transfer 
   TABLE DATA           �   COPY public.tb_group_transfer (id, "createdAt", profile, old_group, new_group, new_group_combination, description, added_tags_record) FROM stdin;
    public          postgres    false    230   ��       [          0    16509 	   tb_groups 
   TABLE DATA           �   COPY public.tb_groups (id, "createdAt", "updatedAt", "deletedAt", name, code, parent, levels, coefficient, locked, invalid, is_default) FROM stdin;
    public          postgres    false    225   j�       r          0    24600 
   tb_message 
   TABLE DATA           v   COPY public.tb_message (id, "createdAt", "updatedAt", "deletedAt", text_id, rec_id, status, message_type) FROM stdin;
    public          postgres    false    248   ��       t          0    24609    tb_message_text 
   TABLE DATA           l   COPY public.tb_message_text (id, send_id, title, text, message_type, "group", "postDate", role) FROM stdin;
    public          postgres    false    250   ��       R          0    16460    tb_operate_record 
   TABLE DATA           B   COPY public.tb_operate_record (id, "createdAt", body) FROM stdin;
    public          postgres    false    216   ��       j          0    16603    tb_permissions 
   TABLE DATA           2   COPY public.tb_permissions (id, name) FROM stdin;
    public          postgres    false    240   ��       b          0    16551 
   tb_profile 
   TABLE DATA           F  COPY public.tb_profile (id, "createdAt", "updatedAt", "deletedAt", name, job_number, type_card, phone, id_card, gender, birth_day, source, school, graduation_date, specialty, last_company, first_job_date, work_age, nation, marital_status, account_location, address, bank_card, on_board_date, freezed, audit_state) FROM stdin;
    public          postgres    false    232   	�       I          0    16409 	   tb_record 
   TABLE DATA           B   COPY public.tb_record (id, "createdAt", body, object) FROM stdin;
    public          postgres    false    207   �       G          0    16397    tb_roles 
   TABLE DATA           S   COPY public.tb_roles (id, "createdAt", "updatedAt", "deletedAt", name) FROM stdin;
    public          postgres    false    205   ��       n          0    16625 	   tb_salary 
   TABLE DATA           k   COPY public.tb_salary (id, "createdAt", template_account, template, year, month, locked, data) FROM stdin;
    public          postgres    false    244   ��       p          0    16636    tb_salary_config 
   TABLE DATA           u   COPY public.tb_salary_config (id, "createdAt", "updatedAt", "deletedAt", base, tax_threshold, reference) FROM stdin;
    public          postgres    false    246   =�       l          0    16614    tb_salary_fields 
   TABLE DATA           �   COPY public.tb_salary_fields (id, "createdAt", profile_id, salary_id, department_group_id, post_group_id, key, name, alias, value, content, should_tax, is_income, is_deduct, year, month, fit_into_year, fit_into_month) FROM stdin;
    public          postgres    false    242   Z�       K          0    16420    tb_salary_profile_config 
   TABLE DATA           �   COPY public.tb_salary_profile_config (id, "createdAt", "updatedAt", "deletedAt", template_field_id, profile_id, operate, value, description) FROM stdin;
    public          postgres    false    209   y�       P          0    16447    tb_tags 
   TABLE DATA           X   COPY public.tb_tags (id, name, coefficient, parent, commensalism_group_ids) FROM stdin;
    public          postgres    false    214   ��       ^          0    16528    tb_template 
   TABLE DATA           �   COPY public.tb_template (id, "createdAt", "updatedAt", "deletedAt", name, type, months, startup, init_data, "order", user_id, audit_state) FROM stdin;
    public          postgres    false    228   �       d          0    16566    tb_template_account 
   TABLE DATA           M   COPY public.tb_template_account (id, "createdAt", name, "order") FROM stdin;
    public          postgres    false    234   �      h          0    16591    tb_usergroups 
   TABLE DATA           h   COPY public.tb_usergroups (id, "createdAt", "updatedAt", "deletedAt", name, parent, levels) FROM stdin;
    public          postgres    false    238   �      U          0    16476    tb_users 
   TABLE DATA           �   COPY public.tb_users (id, "createdAt", "updatedAt", email, username, nichname, id_card, password, is_super, picture, state) FROM stdin;
    public          postgres    false    219   l      Y          0    16502    templateaccount_groups 
   TABLE DATA           O   COPY public.templateaccount_groups (group_id, template_account_id) FROM stdin;
    public          postgres    false    223   �      \          0    16521    templateaccount_templates 
   TABLE DATA           U   COPY public.templateaccount_templates (template_id, template_account_id) FROM stdin;
    public          postgres    false    226   �      V          0    16487    user_groups 
   TABLE DATA           8   COPY public.user_groups (group_id, user_id) FROM stdin;
    public          postgres    false    220   �      D          0    16385 
   user_roles 
   TABLE DATA           6   COPY public.user_roles (role_id, user_id) FROM stdin;
    public          postgres    false    202   	      L          0    16430 	   user_tags 
   TABLE DATA           4   COPY public.user_tags (tag_id, user_id) FROM stdin;
    public          postgres    false    210   *      S          0    16469    user_usergroups 
   TABLE DATA           A   COPY public.user_usergroups (user_id, user_group_id) FROM stdin;
    public          postgres    false    217   G      �           0    0    tb_audit_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.tb_audit_id_seq', 61, true);
          public          postgres    false    235            �           0    0    tb_group_transfer_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.tb_group_transfer_id_seq', 7, true);
          public          postgres    false    229            �           0    0    tb_groups_id_seq    SEQUENCE SET     ?   SELECT pg_catalog.setval('public.tb_groups_id_seq', 69, true);
          public          postgres    false    224            �           0    0    tb_message_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.tb_message_id_seq', 34, true);
          public          postgres    false    247            �           0    0    tb_message_text_id_seq    SEQUENCE SET     E   SELECT pg_catalog.setval('public.tb_message_text_id_seq', 39, true);
          public          postgres    false    249            �           0    0    tb_operate_record_id_seq    SEQUENCE SET     H   SELECT pg_catalog.setval('public.tb_operate_record_id_seq', 188, true);
          public          postgres    false    215            �           0    0    tb_permissions_id_seq    SEQUENCE SET     D   SELECT pg_catalog.setval('public.tb_permissions_id_seq', 1, false);
          public          postgres    false    239            �           0    0    tb_profile_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.tb_profile_id_seq', 57, true);
          public          postgres    false    231            �           0    0    tb_record_id_seq    SEQUENCE SET     ?   SELECT pg_catalog.setval('public.tb_record_id_seq', 61, true);
          public          postgres    false    206            �           0    0    tb_roles_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.tb_roles_id_seq', 6, true);
          public          postgres    false    204            �           0    0    tb_salary_config_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.tb_salary_config_id_seq', 1, false);
          public          postgres    false    245            �           0    0    tb_salary_fields_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.tb_salary_fields_id_seq', 38, true);
          public          postgres    false    241            �           0    0    tb_salary_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.tb_salary_id_seq', 9, true);
          public          postgres    false    243            �           0    0    tb_salary_profile_config_id_seq    SEQUENCE SET     M   SELECT pg_catalog.setval('public.tb_salary_profile_config_id_seq', 1, true);
          public          postgres    false    208            �           0    0    tb_tags_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.tb_tags_id_seq', 12, true);
          public          postgres    false    213            �           0    0    tb_template_account_id_seq    SEQUENCE SET     I   SELECT pg_catalog.setval('public.tb_template_account_id_seq', 13, true);
          public          postgres    false    233            �           0    0    tb_template_id_seq    SEQUENCE SET     A   SELECT pg_catalog.setval('public.tb_template_id_seq', 40, true);
          public          postgres    false    227            �           0    0    tb_usergroups_id_seq    SEQUENCE SET     B   SELECT pg_catalog.setval('public.tb_usergroups_id_seq', 2, true);
          public          postgres    false    237            �           0    0    tb_users_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.tb_users_id_seq', 8, true);
          public          postgres    false    218            �           2606    16444    group_tags group_tags_pkey 
   CONSTRAINT     f   ALTER TABLE ONLY public.group_tags
    ADD CONSTRAINT group_tags_pkey PRIMARY KEY (tag_id, group_id);
 D   ALTER TABLE ONLY public.group_tags DROP CONSTRAINT group_tags_pkey;
       public            postgres    false    212    212            ~           2606    16394 (   permissions_roles permissions_roles_pkey 
   CONSTRAINT     {   ALTER TABLE ONLY public.permissions_roles
    ADD CONSTRAINT permissions_roles_pkey PRIMARY KEY (role_id, permissions_id);
 R   ALTER TABLE ONLY public.permissions_roles DROP CONSTRAINT permissions_roles_pkey;
       public            postgres    false    203    203            �           2606    16496 "   profile_groups profile_groups_pkey 
   CONSTRAINT     r   ALTER TABLE ONLY public.profile_groups
    ADD CONSTRAINT profile_groups_pkey PRIMARY KEY (group_id, profile_id);
 L   ALTER TABLE ONLY public.profile_groups DROP CONSTRAINT profile_groups_pkey;
       public            postgres    false    221    221            �           2606    16439    profile_tags profile_tags_pkey 
   CONSTRAINT     l   ALTER TABLE ONLY public.profile_tags
    ADD CONSTRAINT profile_tags_pkey PRIMARY KEY (tag_id, profile_id);
 H   ALTER TABLE ONLY public.profile_tags DROP CONSTRAINT profile_tags_pkey;
       public            postgres    false    211    211            �           2606    16501    tag_groups tag_groups_pkey 
   CONSTRAINT     f   ALTER TABLE ONLY public.tag_groups
    ADD CONSTRAINT tag_groups_pkey PRIMARY KEY (group_id, tag_id);
 D   ALTER TABLE ONLY public.tag_groups DROP CONSTRAINT tag_groups_pkey;
       public            postgres    false    222    222            �           2606    16587    tb_audit tb_audit_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.tb_audit
    ADD CONSTRAINT tb_audit_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.tb_audit DROP CONSTRAINT tb_audit_pkey;
       public            postgres    false    236            �           2606    16548 (   tb_group_transfer tb_group_transfer_pkey 
   CONSTRAINT     f   ALTER TABLE ONLY public.tb_group_transfer
    ADD CONSTRAINT tb_group_transfer_pkey PRIMARY KEY (id);
 R   ALTER TABLE ONLY public.tb_group_transfer DROP CONSTRAINT tb_group_transfer_pkey;
       public            postgres    false    230            �           2606    16519    tb_groups tb_groups_pkey 
   CONSTRAINT     V   ALTER TABLE ONLY public.tb_groups
    ADD CONSTRAINT tb_groups_pkey PRIMARY KEY (id);
 B   ALTER TABLE ONLY public.tb_groups DROP CONSTRAINT tb_groups_pkey;
       public            postgres    false    225            �           2606    24605    tb_message tb_message_pkey 
   CONSTRAINT     X   ALTER TABLE ONLY public.tb_message
    ADD CONSTRAINT tb_message_pkey PRIMARY KEY (id);
 D   ALTER TABLE ONLY public.tb_message DROP CONSTRAINT tb_message_pkey;
       public            postgres    false    248            �           2606    24617 $   tb_message_text tb_message_text_pkey 
   CONSTRAINT     b   ALTER TABLE ONLY public.tb_message_text
    ADD CONSTRAINT tb_message_text_pkey PRIMARY KEY (id);
 N   ALTER TABLE ONLY public.tb_message_text DROP CONSTRAINT tb_message_text_pkey;
       public            postgres    false    250            �           2606    16468 (   tb_operate_record tb_operate_record_pkey 
   CONSTRAINT     f   ALTER TABLE ONLY public.tb_operate_record
    ADD CONSTRAINT tb_operate_record_pkey PRIMARY KEY (id);
 R   ALTER TABLE ONLY public.tb_operate_record DROP CONSTRAINT tb_operate_record_pkey;
       public            postgres    false    216            �           2606    16611 "   tb_permissions tb_permissions_pkey 
   CONSTRAINT     `   ALTER TABLE ONLY public.tb_permissions
    ADD CONSTRAINT tb_permissions_pkey PRIMARY KEY (id);
 L   ALTER TABLE ONLY public.tb_permissions DROP CONSTRAINT tb_permissions_pkey;
       public            postgres    false    240            �           2606    16562 !   tb_profile tb_profile_id_card_key 
   CONSTRAINT     _   ALTER TABLE ONLY public.tb_profile
    ADD CONSTRAINT tb_profile_id_card_key UNIQUE (id_card);
 K   ALTER TABLE ONLY public.tb_profile DROP CONSTRAINT tb_profile_id_card_key;
       public            postgres    false    232            �           2606    16560    tb_profile tb_profile_pkey 
   CONSTRAINT     X   ALTER TABLE ONLY public.tb_profile
    ADD CONSTRAINT tb_profile_pkey PRIMARY KEY (id);
 D   ALTER TABLE ONLY public.tb_profile DROP CONSTRAINT tb_profile_pkey;
       public            postgres    false    232            �           2606    16417    tb_record tb_record_pkey 
   CONSTRAINT     V   ALTER TABLE ONLY public.tb_record
    ADD CONSTRAINT tb_record_pkey PRIMARY KEY (id);
 B   ALTER TABLE ONLY public.tb_record DROP CONSTRAINT tb_record_pkey;
       public            postgres    false    207            �           2606    16405    tb_roles tb_roles_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.tb_roles
    ADD CONSTRAINT tb_roles_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.tb_roles DROP CONSTRAINT tb_roles_pkey;
       public            postgres    false    205            �           2606    16644 &   tb_salary_config tb_salary_config_pkey 
   CONSTRAINT     d   ALTER TABLE ONLY public.tb_salary_config
    ADD CONSTRAINT tb_salary_config_pkey PRIMARY KEY (id);
 P   ALTER TABLE ONLY public.tb_salary_config DROP CONSTRAINT tb_salary_config_pkey;
       public            postgres    false    246            �           2606    16622 &   tb_salary_fields tb_salary_fields_pkey 
   CONSTRAINT     d   ALTER TABLE ONLY public.tb_salary_fields
    ADD CONSTRAINT tb_salary_fields_pkey PRIMARY KEY (id);
 P   ALTER TABLE ONLY public.tb_salary_fields DROP CONSTRAINT tb_salary_fields_pkey;
       public            postgres    false    242            �           2606    16633    tb_salary tb_salary_pkey 
   CONSTRAINT     V   ALTER TABLE ONLY public.tb_salary
    ADD CONSTRAINT tb_salary_pkey PRIMARY KEY (id);
 B   ALTER TABLE ONLY public.tb_salary DROP CONSTRAINT tb_salary_pkey;
       public            postgres    false    244            �           2606    16428 6   tb_salary_profile_config tb_salary_profile_config_pkey 
   CONSTRAINT     t   ALTER TABLE ONLY public.tb_salary_profile_config
    ADD CONSTRAINT tb_salary_profile_config_pkey PRIMARY KEY (id);
 `   ALTER TABLE ONLY public.tb_salary_profile_config DROP CONSTRAINT tb_salary_profile_config_pkey;
       public            postgres    false    209            �           2606    16457    tb_tags tb_tags_pkey 
   CONSTRAINT     R   ALTER TABLE ONLY public.tb_tags
    ADD CONSTRAINT tb_tags_pkey PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.tb_tags DROP CONSTRAINT tb_tags_pkey;
       public            postgres    false    214            �           2606    16576 0   tb_template_account tb_template_account_name_key 
   CONSTRAINT     k   ALTER TABLE ONLY public.tb_template_account
    ADD CONSTRAINT tb_template_account_name_key UNIQUE (name);
 Z   ALTER TABLE ONLY public.tb_template_account DROP CONSTRAINT tb_template_account_name_key;
       public            postgres    false    234            �           2606    16574 ,   tb_template_account tb_template_account_pkey 
   CONSTRAINT     j   ALTER TABLE ONLY public.tb_template_account
    ADD CONSTRAINT tb_template_account_pkey PRIMARY KEY (id);
 V   ALTER TABLE ONLY public.tb_template_account DROP CONSTRAINT tb_template_account_pkey;
       public            postgres    false    234            �           2606    16536    tb_template tb_template_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.tb_template
    ADD CONSTRAINT tb_template_pkey PRIMARY KEY (id);
 F   ALTER TABLE ONLY public.tb_template DROP CONSTRAINT tb_template_pkey;
       public            postgres    false    228            �           2606    16599     tb_usergroups tb_usergroups_pkey 
   CONSTRAINT     ^   ALTER TABLE ONLY public.tb_usergroups
    ADD CONSTRAINT tb_usergroups_pkey PRIMARY KEY (id);
 J   ALTER TABLE ONLY public.tb_usergroups DROP CONSTRAINT tb_usergroups_pkey;
       public            postgres    false    238            �           2606    16486    tb_users tb_users_id_card_key 
   CONSTRAINT     [   ALTER TABLE ONLY public.tb_users
    ADD CONSTRAINT tb_users_id_card_key UNIQUE (id_card);
 G   ALTER TABLE ONLY public.tb_users DROP CONSTRAINT tb_users_id_card_key;
       public            postgres    false    219            �           2606    16484    tb_users tb_users_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.tb_users
    ADD CONSTRAINT tb_users_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.tb_users DROP CONSTRAINT tb_users_pkey;
       public            postgres    false    219            �           2606    16506 2   templateaccount_groups templateaccount_groups_pkey 
   CONSTRAINT     �   ALTER TABLE ONLY public.templateaccount_groups
    ADD CONSTRAINT templateaccount_groups_pkey PRIMARY KEY (group_id, template_account_id);
 \   ALTER TABLE ONLY public.templateaccount_groups DROP CONSTRAINT templateaccount_groups_pkey;
       public            postgres    false    223    223            �           2606    16525 8   templateaccount_templates templateaccount_templates_pkey 
   CONSTRAINT     �   ALTER TABLE ONLY public.templateaccount_templates
    ADD CONSTRAINT templateaccount_templates_pkey PRIMARY KEY (template_id, template_account_id);
 b   ALTER TABLE ONLY public.templateaccount_templates DROP CONSTRAINT templateaccount_templates_pkey;
       public            postgres    false    226    226            �           2606    16491    user_groups user_groups_pkey 
   CONSTRAINT     i   ALTER TABLE ONLY public.user_groups
    ADD CONSTRAINT user_groups_pkey PRIMARY KEY (group_id, user_id);
 F   ALTER TABLE ONLY public.user_groups DROP CONSTRAINT user_groups_pkey;
       public            postgres    false    220    220            |           2606    16389    user_roles user_roles_pkey 
   CONSTRAINT     f   ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (role_id, user_id);
 D   ALTER TABLE ONLY public.user_roles DROP CONSTRAINT user_roles_pkey;
       public            postgres    false    202    202            �           2606    16434    user_tags user_tags_pkey 
   CONSTRAINT     c   ALTER TABLE ONLY public.user_tags
    ADD CONSTRAINT user_tags_pkey PRIMARY KEY (tag_id, user_id);
 B   ALTER TABLE ONLY public.user_tags DROP CONSTRAINT user_tags_pkey;
       public            postgres    false    210    210            �           2606    16473 $   user_usergroups user_usergroups_pkey 
   CONSTRAINT     v   ALTER TABLE ONLY public.user_usergroups
    ADD CONSTRAINT user_usergroups_pkey PRIMARY KEY (user_id, user_group_id);
 N   ALTER TABLE ONLY public.user_usergroups DROP CONSTRAINT user_usergroups_pkey;
       public            postgres    false    217    217            �           1259    16588    idx_tb_audit_deletedat    INDEX     R   CREATE INDEX idx_tb_audit_deletedat ON public.tb_audit USING btree ("deletedAt");
 *   DROP INDEX public.idx_tb_audit_deletedat;
       public            postgres    false    236            �           1259    16520    idx_tb_groups_deletedat    INDEX     T   CREATE INDEX idx_tb_groups_deletedat ON public.tb_groups USING btree ("deletedAt");
 +   DROP INDEX public.idx_tb_groups_deletedat;
       public            postgres    false    225            �           1259    24606    idx_tb_message_deletedat    INDEX     V   CREATE INDEX idx_tb_message_deletedat ON public.tb_message USING btree ("deletedAt");
 ,   DROP INDEX public.idx_tb_message_deletedat;
       public            postgres    false    248            �           1259    16563    idx_tb_profile_deletedat    INDEX     V   CREATE INDEX idx_tb_profile_deletedat ON public.tb_profile USING btree ("deletedAt");
 ,   DROP INDEX public.idx_tb_profile_deletedat;
       public            postgres    false    232                       1259    16406    idx_tb_roles_deletedat    INDEX     R   CREATE INDEX idx_tb_roles_deletedat ON public.tb_roles USING btree ("deletedAt");
 *   DROP INDEX public.idx_tb_roles_deletedat;
       public            postgres    false    205            �           1259    16645    idx_tb_salary_config_deletedat    INDEX     b   CREATE INDEX idx_tb_salary_config_deletedat ON public.tb_salary_config USING btree ("deletedAt");
 2   DROP INDEX public.idx_tb_salary_config_deletedat;
       public            postgres    false    246            �           1259    16429 &   idx_tb_salary_profile_config_deletedat    INDEX     r   CREATE INDEX idx_tb_salary_profile_config_deletedat ON public.tb_salary_profile_config USING btree ("deletedAt");
 :   DROP INDEX public.idx_tb_salary_profile_config_deletedat;
       public            postgres    false    209            �           1259    16537    idx_tb_template_deletedat    INDEX     X   CREATE INDEX idx_tb_template_deletedat ON public.tb_template USING btree ("deletedAt");
 -   DROP INDEX public.idx_tb_template_deletedat;
       public            postgres    false    228            �           1259    16600    idx_tb_usergroups_deletedat    INDEX     \   CREATE INDEX idx_tb_usergroups_deletedat ON public.tb_usergroups USING btree ("deletedAt");
 /   DROP INDEX public.idx_tb_usergroups_deletedat;
       public            postgres    false    238            N      x������ � �      E      x������ � �      W   n   x�5���0�f1��:l����X�F@>&�s:�!��&s� �'�8Do4�)6Ŧؔ�����Mf_��A��Kn#��%7��]�򐇶.]ך�ͤfR����D)"�      M       x�3�45�2�45� �`�D��qqq Km      X      x�33�4����� �]      f   f	  x��ZKo�>���=&�=������K�CdX�(�a%�� ��6C=HEѓ�%�V$�
E�,�R�gv��)!U3�ݝ��&�Z\��Ȫ���tC����h ��2��S�g�פ.G�+�z�A����=���G���ӧ��ޟ�/ջ�������y��Q��쟻P�p���hse����_�4(.�`A��J�~1}w����?7��w���|�X�a�޹�BA)�9�K�ѵ��������)X;P*W�x3e��T�JiI�j4�Z���
N��6���N*��xTJ͕��է����z�p�Bn\�Cigf�L���+�@-l�b#�ӝ��|�������O�Շڙ���`����7��be;�{py�z��X���������g`Q���M���N�ϙn���Q���۝K��F1EHL	P�.��X���e�@
Z�fʳV�A�2}T26"S��7|��C���N��p)G�9�N�2�k�y���v���9*���\)��>��j�:��ׁ�5M���k��3|v�XݬX��*��*����[snaŌ���X�x�֟k7����w�(�Y2)V3�� ԘjkC-�M�d3�
�#u?G�6�x�����8�E%�6�Jy21t�C��d��Lc�0�pL聡���i��9�з`g)c��+�ZvRs`'�������1�N����D����-(J�~ZP���H�a2
�ô�0������$|Ĝ�ڥ���������Ý�NMG�H��I���&1�fF�j�	/���5���KKɮ�����!5-����m�)�^{���T�M�������9��O/��r;�5�ۍn���8�S��F�|����X8`@T3is�r�g��X�s$�,� ��MK0��J�fL �5@��Ȁ,x�
H�X��q�Y�d���Tv$W�[n�ul���ɜS�0eҤL6�r{�3����ua�A�Ȫ�`�u�2^��X��]&m�$�H�ѲpAij�	���A�ur*�<3�PO���.�>eRה�\i.y�U�R�uH���(�h��P	�A�x�$.n�b�V
OF��@�fH���Z��"^bҳVWg�J�k��'��	)�i���x�ӎZl��N���MkE�.x޹�(��qg���J>�b��D&36hZl:)�|�֊0���>�����3= l��3"�\2P�R6C������	M��B�O��	�ZմrR�OM*��Y��.eb�e�9e�NI��S��b�&'��A��&FZe���h۴�IHW7�I�L���Z��e\3%H��P�暻�MT)����{���6�6W�
�8��6SiHSS+*��7���Q3'#S2T;y;�֭b�N��������V�v��p�X��_����`V6�7~"=�OO����,��Ι���}���_I��<����e�>��Ifg[��G�!�d�J�Y�1�������?>[��7���p��k�������� ���t��PvY�&S+v��4�|t����m�*֮���w>o���2̝͞
9cZp��z2�Ε؝+�ۨ&2)���2��93B 
���η�12�a�>#o�x)5����Z�o7F��������Hy~��3٭`w�_l�R�8��V�Z���6����4/���q��L��F�ߖ�*�{��� �#�20Dh��J�eF��6H��m��� �{��t�R22����u�����;FFV���o��ʹ���q�N\T�y**{3���^�\^6�q'��J����2�7]�H�����U�C.�`2�N�9���鵶�7�G�I��Z�4�S)�S%��>p���i2��,����Z�3.q/��q�˜JS�ǥvdo��'�@���S�$⹷��.�F�ʊԦ���H��
�Y�|d�YI��!x�
)��
�.�G0I#��;�g�g�.��KF��.��7ޓ�ܕ�_���[�$��W_l>����ꓟ�n?}��������+�ʯŅ���O��6���-3 ��E��-�YZS)q���D�<�yN��y�0|�>�\O<( K=�BA�N�;.2�$ R&`���#:=�r����ѣ��	��ȿ��.��ϼ�~��dy�CR2��m��c����P��L�wQ0�m�$���z�����W[��w�{C����?��r\��1�y1 � NM��&t1�y6!�KP14( ʇ��R���O}���ZYv����<+�]���@8]�8��W�v�����^�y���7��?+�����K'*]'���c�\k�s�X��]�9_EQ�e�x*1����v�ha����O�[�a׷j�O�H�)��t�׭���1�,��ҝ2�<�_�Vo��[-���īkp�`y���7��v�Z��Hy#�	d�6����Ҹ^����u蹱�r?��f�~��!k�      `   �   x�e���0E���
��|ތ?EP�D�=1������+��������ҭ�h/�$�$B�F��$=�-���V� 1t!�zߒ=5,s1u:.w�]#G�a��:�?�q������SchC��� �t����ˍC:�\���n�W��R��a;�      [   .  x���_N�@Ɵ�S�^��ٙ���=�9�TU�mT���Ԋ����R*�6qn�Yۀg�,�������7k.z,��]G^C.ɖl���Sxϩ/{w�&wǓ���k�7��3���ZnuqV��~�u�;)Y�RK̽��������B0z����]��7j�W�Vw;f�Yv�
{��p��j~�I�;"��7j������Be{�]��4��<�~����.�,�S��|�U��	(�l�Z5�/$�*[�6! fѵ�ȿ����&���&��7�p[���e��ǚ,g�ת)z�W��V�K������Z5q�n�dJ�`<>`���bŞ�j��p8�il��m�e����q1���7;�1�����{}͉,:�׌%&0!�{X�NԞ���R�a��̆7�/����韝jt"d�>�=���x��ʲ7�,���ƙe��/W�--������TK
\���kù<+tiD�;�N~V��(]�ꈎ�:����y�9Pa21vn�ժ)���;YmZqd{Y�&e&�3�/|4*UL��V#�}�aiXE�0�~�xJ�M      r   �   x�u�KnB1��{�`^a]�ǋ@] #U��o� �-�We�У���&�J&�_�iZ�9���C�6=�6��{}�]~��b�����3�*bx:y���kYy�Z)
K��ٳ'���~A���ߚ@2�X)!�����>�(�kO�Rlx<y�C~���fkÓ)?�6<��J���v��#�2y��q�ik���ܢ�1���i���#~�      t   �   x��ѽJ�@�z��e�{�ܙ;3O�جtA��7X�B�E���tk��4���
n i�7p��|�8��(+U���|�i����(��q���OgG��̄ G�#�&���^Afp7�)��VP��%�]
�8 J/�1��l� �I �$;F�N���6R�p�,��n�V���y\�������u�?��kz�����(Z�Q��p�tc�:�ܦ��lq(�P�3.�����-���8�EBp�فβ����      R   E  x���Ϗ\������"%��tuU�zRnHW�[��J6Ē�-�-q�`b��(@0؋�I�a���3;������~�;�of���`�+�|^���oUW<c��'�9�iРkЪ�Z��Ϭ�|���������o~�x���c�|$6dZ"�<��#�6@Jpq�#����i�6V+����kF0zuti�ԅ�+/�jF�C^�;���E���/Ǐ���p��ݍ�&��?>����=�h4�������?�����o6:(�#z�E�z �{k�;�?�ܨ�~���_\:}���_�/G�q{7ҺцIJ[0+1�Vddǲ�M����=jL�4���3�g�qaf�ccxq�"�c���i�����׏���񍷚��g�N�_޾�.�X�"��d�
�x���G�6~4����@�<ߐV����g�$m}ye|}��볿%��
d��X�T�߾ӄi�`o��CQۣ�����[��^��q�S4�1~rg�ᛉAsAYcL����]�5lJ4��J�m�S�uG���o������^��.�
���f�ƃ�;7�O?O�0C*�Cd�\�8��E�p�M��L� cZ8�)�V�0�-d�p���e4�[k�����ջ9&�o�����oST�t�̹��Ϟ~������9�N��z~����6@@:0��3���0bL�nKʓ�~�9�y�� �/T��Β⿡.�m#�8�.���%cT r�޷�qn�)��v�`\�7Ɨ ^$�Ѓ�����'�o��u|���Gak/��E�mH���"0�إ�j�`B��<`���B]�D�z� ���,}c�(�d�Z?�oOV8�0l����]���[� � D圣~�b�r��/[��O�R81�G֛S[`l�P���j��Ɲ��\4Z�n�5�
ؐ�#��B�����
���8��:����)"s��<+��[2+w�m~����++��~8y�ﭵ�[*vQ�"�"�Z��u���{��|��]���E�>&m�[t!����:�F�E���_����@�v~L[��(���զL� [O��3#���B�D��"Ȋ|>���+L+�ŦZ[.����İ]�b^�G���Փ/������9 ���eb��7]+��_�z��Z�A��%�ʹ�
Йj�sl(���VT�!����Z��֮���om>�7Y_�|����+�X�9�|���o^��ҕ��*:�չ��<g�f@��:pn�]�.o{�e�䀥=A�e��K`h4[�9xH(�d�RB���t�`���j�\�1��GR�h�g�0'���\��n[�K�Q!�;]�'o~�}흜Qf�N���D}��+)��o.��t��S���\�I���/^�g%��*2����x���������������K�/�:�r�����<�y�P�fL�������}�X���R�'}�^�[:]��UX�?�kB�h��j�:��o9�CsV!nT���Cv�%��Ƈ�ҟ��D�F���HX;�,�] �?���� ���\	�,_-<�Y�K`�D���R4��P�2�X��K �A�)1�����e%aTpa)z�zI�������}(	��T;�Ns��HbACN��O����ꛟ�'.���|�`���^���3X�tR�a7=R�C`]�e1��7�cuM�X_b�X�1��Z�a���Zl)
�^Q�K*I$�ա�ؠK�������%��J,����r���,k ɀPv�+*`�ur`�X�e�-�X/�ǀ6���X�*�� =Jm���.���On�.9�9X��^�Q,Q�k��Zc'��n��4�&�}+�ŕ��?#?�|:�Ս��D�Z���߾�������.>�tc�����o��� iL�eO���K�D��l@�HU]�_v�Z�$�~����2e?˨�S9�I��V�|��Q
���7��?�9%���w�0�Q��5��4?��1�t�jrNC	�-(���5�X`Z�0�h��� ��br���wX-+Ex5ߝo�ҟ�Y�(��o\l��u$I���X]��kv�y���uk�$����t��"�l�$�ѣun��Bk�+H`���V��^�q�D���	u��im]K�;U7���-�H�(��@�3l\AwI6(5-kU�H��y����0�cض�B�O�ĉ{ЄՇ3}�+��4��c�A4sm,��DHQ"����& �?1Z�sҪ:��öUJ"$/=�aI!`֟�J	��`���}��D�T��VǁƳ�J�|��v����E�uC��8� �Z�`�K���m�g�H���i�0��l0X���l�����q�E)�؉#�����6��,ԏZ�_d���<�R���S:�Ԥ�{����׸ҡ0��y9hڹS:")��1��P�t($2S9�����ku�E���k����t(�����3�lke���Ar(ZR��ҝP*��Xݷw�m�j��;	ҳ>j��@��;�N.Z��M��_��s!�6W�p֝��Ӧ�t��q����R�Z�`\�H����9��ꃮ��t�@j0�l[�@3�� ��Իu���ߒpɧ;XCI�Ҕ	#��B�С��%�w�l�$-'9�@����	ku�t��	�rrAsI�t��8�Nc���X�Su��Ƶ:h��ӵ�t�m����[z������;�[b�=�m��Q���Qf9��K<���z��=` [���L�I)�q
��-者݆`������_�Y��+`��7?�v�.uԹ�q�C��PJi_EN� ,|Ȗ�t!�E�Go�C�x�R	��a�$�n[l~�X.%VVKl�����:�XP��P�#Sy{q����	�@���cu������N$i<�!o�ʍ�k�'���0jE�Ɓ�ݹan
�|C{d�u��yK���C�\B�Z뫏k�Z;�i�����At�e��|Ћ��2&ڒ�Z����	���V�;�%QT1�i��a� �4��{��i[�]��q��q� e�IC�B���C�ZSv�}��yUߌb�1Ik��9Ԡ��X�����Y���$B�(��F.o���ܜ�g�GՕ�>�bL;,S[wB�N�}�E��x�!�tQ�1T��-�����
����Cq��J� �@���MR�B�#Gy�vwF�R8D��	�7��5���M�W_\]9�	�$�p�\uY�x�����HQMpֆ3�ξ���D��-���eyx�Ͼ���S[;��dS�mٍ��Io1���)�I��4E���8`�PjPF�^�-_��#�_(d�����iV2`��k{W}��>>忦�t3ܔ�T�J�0هk�	ao�އ(C��o�?g�Ա�'ACJr����Q]N���;�8~��      j      x������ � �      b   �  x����jA��5O�{�Υn�n\��&A\�{W*Y��QA#��"T�i���-<�kL�{&(#S�������N{0���k�
�!jm9��&"C㒍����;�L��u;?<�yd��ۇ��?����J9f���sr�̯�'FjkȦ{s���4��{��o�U{�wo����E��g�M��}���bX+807�<��RP����`�� �2���~=�6��� `��g��k'Vh{���	M�Bnl����I����=l,��oG$H�CL�H��&\�A�@�"�_$R�a{���l3��<٨j$_����j=�)	��t۷T	����#�bR(���
O����s�	�H�%���,�k ����O4�S�3�D�􉼖�?$^��څ�����K/IGYf��p����'�ۓ�^��f:�A��t^��cVM����%����)[y�0�b��.�*)�x���W�����^w�l6��n�@      I   �  x��Z]OSI��_��ݬ�3�|O�w{��`��u7�%�hL�"j�SWE\�B`�lD��g:�����}���)� Z��$����y�wޏ��>�@/R�He�i�QCB~GM_45U����M������\����7����?��oo6�w�/?�|z�-���~�Ciq��C ��Rͬ�Le����A�Z'��0���AL�ʛw���I���I4�f��K������I�͌+�<�F��l=��C{~jƏO_�$f4��:��y���0�@:���kթF~o�����T�D���J�x�O8f�DTw�"�L�q�(�Y�|@�!�E8WTw<�cB��`�� �=�h��پ�����������a�Q�(%}����Ų���$,|n��b��{�XX.n�%߹�{��H߳�W�f��������V6��5�f��߸���W�mzd
G����Y�=��}h���Q��9|��4{mM�>J����	���j�aUi0cPǥ�3��`V�M��R�0Rܞ����T���ޖ�����H�Ds�da$�l�PΏ�Q�5��d�䂿�K�2y��Ԏ�r���K���á���WK�o7�ݺ�Ϫ0Ԅl�c��6��m�4�j>z��=�uf*�aܷ�p�¬�.@4ƃ:BS%Y3Yg�B�P���f��*vS�%a�PhF��O�z	�:m0Ũ���F�g|;ٸ�7����O���.�>�H᎘To*��'���hu�Lz�����]���R;m��I��%��%��B�W��h����?Wg���t�X�c�)Ǩ�T�.W������|$��!=B�y�h*(����=O�ǒ�eC��}�]&/��c��%�"L(�ٹH�~}U��4yv�1J�|e���J��!� ���\�#��q��y�i´� zLޝZy�H�3����k�V�A^I1T����+d�<��k�����L�io��q�2�%�2*u3H�鰵�@��g9O�a8"�-�D�L�N%ݢ��,���#�`-�`�RV�sSI,�Z��6��p#�i��ԙ)�u"5cZ63�s�~}��^��1t�hb�p$��'ŏ��񌟚���Х�MAs�����������B������,�qi)��. �k��A6]�]B�id�!�qM�.�.A�4�{��K��n�u	Z��URBI�|��KЪ:��nY��u
Z��k�|�5r(ǊۣŽ�h�A���}2η��^��˥�W�ܸs��:�s�ՙ�Jhr{���ˁ\}��V$���)���c]lZ�[�L0y�)9�0(���ȃ���ڸ�`� h:m�dTz���;��������n��θ���wJs�w���W��^�`eQ�?��������z���<0cUf6>��NaA�^���+�V����:/%��X8��K�Pey�#��Σ(���BTajS�S�a�!���v
Yy��6�7��C2�ߨ������8��,Bo��:cW{)��!��&V��ϣߍ�7�e�3/�W�H��`���o���k8A��|.����\��L�;j�@��im����&�E�\b����ϯ�bar y�e���n{t�7�`u 8Vk�Ǝ��;̡�:@L�ԑa�{6f46�_txj������цuu	�#T6�S0��N�9^9��hs�s8S�ZH����pEc��o�ۇ�PR}Xn=]
�?��V	��CJ�b-�eg�,�\m]DH��������?B      G   �   x���1
�0��9�������i�	z�
���Z7)��i��l]��_~��#�0�p��%HB)>����̻祻��¬��O�T�i�U��y��}Ä�t��j)s*�C�(
�3�:L�RM��t��S)9��Z�{+K)�X�TAi��R�s�T$����u�      n   }   x�3�4202�50�52W04�26�20�32�4��6��4�|:׳9k�n_�bk˳�����i`ę���e5�X��d�)��YZBL0�i�1�Kt#,��,�͌M�a1"F��� �BO      p      x������ � �      l     x����JA��3�R�9�o�%zA$����@�(��@0)(��6�i�u}��e���3�����,h����z��Z�Gz��{��a��.���܄�� �2��<l��!O��?>�N����;ѹP�6DRz��wNژw��4<��$�X= "0����k}�W�d�RTqk��])%/�_)X��UO��Yq1Ka����1w�.�E�-����>H��wt��m]O�sv׳���2~���6~��oM�?�_��/����=V      K   X   x��̻�0�ڙ�>2zNBl3����#��'�PB#q�&�;�������������K��MGa��ݎ����U�p�!��M?      P   �   x�3�|�zƳ+_,\�b�NCSN��Z.#Η�^�[�b�NN#��1H��u9-�"F �';�>���i1�|�w�4NC�.3�DNC 2�9Ӏ ε�L,NIK" D0qKΧK7?ݾ�i���S惜 2*Ə�h	Y�y���;�=[��|�>�
#Χk'@ŀ
�j��1z\\\ �T�      ^   �   x����JA��ݧHn3��w{�>A��:�$������B���@����@�eꏙ�A`�8@L��R��*J�#E�/_v�}~|�v�ߟ��)o���dzݴ���ݶ��r4k�fz?���e��y��w�� +��ܼ����Z_�Kp%�Z⑖H��2ҟS���@H��&DD=X�G�c%v쬵?
$      d   G   x�< ��1	2020-02-20 11:17:34.701019+08	基本工资账套	{2}
\.


}�      h   �   x�3�4202�50�54R04�26�22г07701�6�  ���lچ�3�=[��Ӏ�@�����D��������D�������<|�ƺ�f
�V&FV�z�F��`�g�k���m�өm/'o{��bg� ��2�      U   -  x����N�`�u�,��{���+E�)�h�B��IZj���1qJܜ��,Η+5j@k(*�Cq!����?�g��U�W����q7��y��s�n՘�	�y�/a'��s��ϧ���^#�^4�_k̢�`�Q&���h����'Q) 	��y����aO~����ݤhw���l��t��fsd$���:��ݞW��g^�IXѢ��n�~uKf��E�\���e���4Yr��;QI��
�x6����NGl1mͣ ��)�<t�ၶ-X���e��m��'C��OAr�      Y      x�3��4����� �Z      \      x�3�4����� f      V      x������ � �      D      x�3�4����� ]      L      x������ � �      S      x�3�4�2bs � �=... "�     