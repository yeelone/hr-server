PGDMP                         x            db_hr    12.1 (Debian 12.1-1.pgdg90+1)    12.1 (Debian 12.1-1.pgdg90+1) �    w           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
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
    public          postgres    false    212   H�       E          0    16390    permissions_roles 
   TABLE DATA           D   COPY public.permissions_roles (role_id, permissions_id) FROM stdin;
    public          postgres    false    203   e�       W          0    16492    profile_groups 
   TABLE DATA           >   COPY public.profile_groups (group_id, profile_id) FROM stdin;
    public          postgres    false    221   ��       M          0    16435    profile_tags 
   TABLE DATA           :   COPY public.profile_tags (tag_id, profile_id) FROM stdin;
    public          postgres    false    211   �       X          0    16497 
   tag_groups 
   TABLE DATA           6   COPY public.tag_groups (group_id, tag_id) FROM stdin;
    public          postgres    false    222   1�       f          0    16579    tb_audit 
   TABLE DATA           �   COPY public.tb_audit (id, "createdAt", "updatedAt", "deletedAt", operator_id, auditor_id, object, action, fields, org_object_id, dest_object_id, state, reply, body, remark) FROM stdin;
    public          postgres    false    236   S�       `          0    16540    tb_group_transfer 
   TABLE DATA           �   COPY public.tb_group_transfer (id, "createdAt", profile, old_group, new_group, new_group_combination, description, added_tags_record) FROM stdin;
    public          postgres    false    230   h�       [          0    16509 	   tb_groups 
   TABLE DATA           �   COPY public.tb_groups (id, "createdAt", "updatedAt", "deletedAt", name, code, parent, levels, coefficient, locked, invalid, is_default) FROM stdin;
    public          postgres    false    225   :�       r          0    24600 
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
    public          postgres    false    232   �       I          0    16409 	   tb_record 
   TABLE DATA           B   COPY public.tb_record (id, "createdAt", body, object) FROM stdin;
    public          postgres    false    207   �       G          0    16397    tb_roles 
   TABLE DATA           S   COPY public.tb_roles (id, "createdAt", "updatedAt", "deletedAt", name) FROM stdin;
    public          postgres    false    205   S      n          0    16625 	   tb_salary 
   TABLE DATA           k   COPY public.tb_salary (id, "createdAt", template_account, template, year, month, locked, data) FROM stdin;
    public          postgres    false    244         p          0    16636    tb_salary_config 
   TABLE DATA           u   COPY public.tb_salary_config (id, "createdAt", "updatedAt", "deletedAt", base, tax_threshold, reference) FROM stdin;
    public          postgres    false    246   �      l          0    16614    tb_salary_fields 
   TABLE DATA           �   COPY public.tb_salary_fields (id, "createdAt", profile_id, salary_id, department_group_id, post_group_id, key, name, alias, value, content, should_tax, is_income, is_deduct, year, month, fit_into_year, fit_into_month) FROM stdin;
    public          postgres    false    242   �      K          0    16420    tb_salary_profile_config 
   TABLE DATA           �   COPY public.tb_salary_profile_config (id, "createdAt", "updatedAt", "deletedAt", template_field_id, profile_id, operate, value, description) FROM stdin;
    public          postgres    false    209   �      P          0    16447    tb_tags 
   TABLE DATA           X   COPY public.tb_tags (id, name, coefficient, parent, commensalism_group_ids) FROM stdin;
    public          postgres    false    214   "      ^          0    16528    tb_template 
   TABLE DATA           �   COPY public.tb_template (id, "createdAt", "updatedAt", "deletedAt", name, type, months, startup, init_data, "order", user_id, audit_state) FROM stdin;
    public          postgres    false    228   �      d          0    16566    tb_template_account 
   TABLE DATA           M   COPY public.tb_template_account (id, "createdAt", name, "order") FROM stdin;
    public          postgres    false    234   �      h          0    16591    tb_usergroups 
   TABLE DATA           h   COPY public.tb_usergroups (id, "createdAt", "updatedAt", "deletedAt", name, parent, levels) FROM stdin;
    public          postgres    false    238   C      U          0    16476    tb_users 
   TABLE DATA           �   COPY public.tb_users (id, "createdAt", "updatedAt", email, username, nichname, id_card, password, is_super, picture, state) FROM stdin;
    public          postgres    false    219   �      Y          0    16502    templateaccount_groups 
   TABLE DATA           O   COPY public.templateaccount_groups (group_id, template_account_id) FROM stdin;
    public          postgres    false    223   3      \          0    16521    templateaccount_templates 
   TABLE DATA           U   COPY public.templateaccount_templates (template_id, template_account_id) FROM stdin;
    public          postgres    false    226   U      V          0    16487    user_groups 
   TABLE DATA           8   COPY public.user_groups (group_id, user_id) FROM stdin;
    public          postgres    false    220   y      D          0    16385 
   user_roles 
   TABLE DATA           6   COPY public.user_roles (role_id, user_id) FROM stdin;
    public          postgres    false    202   �      L          0    16430 	   user_tags 
   TABLE DATA           4   COPY public.user_tags (tag_id, user_id) FROM stdin;
    public          postgres    false    210   �      S          0    16469    user_usergroups 
   TABLE DATA           A   COPY public.user_usergroups (user_id, user_group_id) FROM stdin;
    public          postgres    false    217   �      �           0    0    tb_audit_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.tb_audit_id_seq', 65, true);
          public          postgres    false    235            �           0    0    tb_group_transfer_id_seq    SEQUENCE SET     G   SELECT pg_catalog.setval('public.tb_group_transfer_id_seq', 10, true);
          public          postgres    false    229            �           0    0    tb_groups_id_seq    SEQUENCE SET     ?   SELECT pg_catalog.setval('public.tb_groups_id_seq', 70, true);
          public          postgres    false    224            �           0    0    tb_message_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.tb_message_id_seq', 41, true);
          public          postgres    false    247            �           0    0    tb_message_text_id_seq    SEQUENCE SET     E   SELECT pg_catalog.setval('public.tb_message_text_id_seq', 39, true);
          public          postgres    false    249            �           0    0    tb_operate_record_id_seq    SEQUENCE SET     H   SELECT pg_catalog.setval('public.tb_operate_record_id_seq', 336, true);
          public          postgres    false    215            �           0    0    tb_permissions_id_seq    SEQUENCE SET     D   SELECT pg_catalog.setval('public.tb_permissions_id_seq', 1, false);
          public          postgres    false    239            �           0    0    tb_profile_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.tb_profile_id_seq', 57, true);
          public          postgres    false    231            �           0    0    tb_record_id_seq    SEQUENCE SET     ?   SELECT pg_catalog.setval('public.tb_record_id_seq', 64, true);
          public          postgres    false    206            �           0    0    tb_roles_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.tb_roles_id_seq', 6, true);
          public          postgres    false    204            �           0    0    tb_salary_config_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.tb_salary_config_id_seq', 1, false);
          public          postgres    false    245            �           0    0    tb_salary_fields_id_seq    SEQUENCE SET     G   SELECT pg_catalog.setval('public.tb_salary_fields_id_seq', 893, true);
          public          postgres    false    241            �           0    0    tb_salary_id_seq    SEQUENCE SET     ?   SELECT pg_catalog.setval('public.tb_salary_id_seq', 41, true);
          public          postgres    false    243            �           0    0    tb_salary_profile_config_id_seq    SEQUENCE SET     M   SELECT pg_catalog.setval('public.tb_salary_profile_config_id_seq', 1, true);
          public          postgres    false    208            �           0    0    tb_tags_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.tb_tags_id_seq', 12, true);
          public          postgres    false    213            �           0    0    tb_template_account_id_seq    SEQUENCE SET     I   SELECT pg_catalog.setval('public.tb_template_account_id_seq', 13, true);
          public          postgres    false    233            �           0    0    tb_template_id_seq    SEQUENCE SET     A   SELECT pg_catalog.setval('public.tb_template_id_seq', 41, true);
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
       public            postgres    false    238            N      x������ � �      E      x������ � �      W   o   x�5��1��`V�lr���X�G�T�|�0d���I؋C$l� =N$q����R,�R,�R�9b�C|��ٯ��A��{��~ɷ|�]��7K��WS5�G~p�w��� �A�"�      M       x�3�45�2�45� �`�D��qqq Km      X      x�33�4����� �]      f   
  x��ZKo�>���=&�=������K�CdX�(�a%�� ��6C=HYѓ�%�V$�
E�,�R�gv��)!U3�ݝ��&�Z\��ȪgwC����h ��2��S��פ.G�+�z�A����=���G���S'���_�oջ�����˭�X^/v��Ξ/~�7��f���T}���rX������X�����Sܻ���K{��w^�=:[���{�ln!�����إ���V��=�\����(�+d��2�X�Q	��$j5t-TpB'TF��f'�Jp4*��J��ʓљ�}=�8P!7.�!��3�l&RC̕L�6T�����Tx�?}����֧�ʃ��p�b��s}o��J�R����޿4Z9�o,��f�^��hr�3��}��&R�N���L7����\����b}���"$��Y|I�A�Z�� �d3�Y�� c��>,�)��>�ǁd|h'cc����Z�]ʵ�<YM]+CL��Ed�o,�Yno�z����@�ך�|�O�=��>�Y�lT�\n�yc��vV��-�9��bF��p�L�Dk���ŋ�ݫd�W�,���f j�
�����&R��rZݑ���#w4�����C�ᢒk�y�<���!�d2�^��z�z8"���AE�4p�͜U�[����k�t=vRs`'�������1�N����D��I[P�<�)���������d��i�a�) Ju�_�!I��9�Keޣ3;{狕����X�8����拓w_Mb,̌&��g^��5���KKɮ�����e!5-����m�)�^{���T�M�������;��O��s�v*�k*�=��~9�S��mF�|�~��X8`@T3is�r�g��X�s$�,� ��MK0��J�fL �5@��Ȁ,x�
H�X��q�Y�d���TV$W�[n�ul���ɜS�0eҤL6�r{�3����ua�A�Ȫ�`�u�2^��X��]&m�$�H��������1vt)���T�yf0��
ۥ]&}ʤ�)u��\�2�����x+hiQ��&�9���� �I.\���&�.��87��͐29��\YE�Ĥg���򕔟�3�t�&�lR�Uҿ�O;j�K;mbj5m�%�p�y�Zl�dg ǝY�f*�������ؠuh��B�Z+��f*�`C+���8��f�Έ�/r�@J��VR�l5�Mh�$2`�|ҤM�Ԫ���:��xB@h�P)%β@���t)�k-��);v�Hj��j�69	�1S�i�	�{v�m�&'!E\�l�&i3��nj�F�q͔ SB��k�S6Q�x��V�U�[o�l�\�S Ǒ���JC��ZQɨF��T�j�9�����ɳX�.6o˷���FW�F6��K�����u����7�����3�q~�/����{�>q��G��3�WҾ!�耸�w�ǧN|������Qq�4١v�sL��>"����OW��~1�y��d���y�$2@�m#��4�]����Ԋ�i.�.]�0|qK��ի�'���ϛ� �C�L�'s�BΘ�"�vFӹ�s�bp�D&TP�7gFD!�����pC��q�s���1^J���V��ƍ�q��t�-~"R����Lv+؝�[����+Ψ���$�&=�u7��mbd\~(�M2�H���2BY�po�]~��bdV��R]I��(>�i|Ӹm��xOԒNXJFF���]'��n|A�cdd娍���&n��L;P�Q[���E5���r�7ӊj���e��w����ʸ����~ӕ�T���P5:�R
&�>�t���Z�^kIxQxԝ49��Is<��?�Q�����.9�&�͂Z�����7���W�8Ꭹ��y�45}Tj����~��/=O"�{��RjT��Hmz{O�f@¦V����#� {%9R����UH�|U�uy=�I�w^I={�8+u��`\2�Lw!������<���,αe'YH��j����f�Ǔ��rt���/�k�?/'����=����ݼ&��8�\.�G��fiM��M����d�9�������hc-y�, ,��"Q��:�$w\d�I@�L���Gtz����G��_w ����_�8���}�S��I��c�u��(���\ə��`��,I
���s{���������c��w�_�"�&F�c�� b ��lM��&t1�y!�KP14( ʇ��R���O}���ZYv����<+�]���@8}�qXw?��e"<���rnȫ��Q|�Y��,�\<V�:^,o��ZÜC�Jvw��|E�����h�Rۅ��qF\kl�F<^.9�]ߪ	?9�#)�2��^���c@.�\���Kwʰ�xk|��BX��>o��\u ^]��0�K�%�� ��k�o"�&�i[Ȓ���J�z��w�ׁ�Ɔ���[�`�J0���N�Z�UB��MPm�w_�ȚÜ�>|8�>� ��V��NS�Z{�B���
T�K�/�?�9|��wrg��G�rct���������2e
s�$YPΡ}[��q���%�!I��"�r��d䴩�������S�U?�|��5�Sޜ�P��������54�0      `   �   x�e�ˍ�0E�5UE���)��T�N��>��B`g.�HF�y#~`v�.Yӛj>��/����W˃��vi�T�}� \z��*C�xj�h�k��ճ�Ll����k�;�^��N-�)���G�xܸR�V�UNn��+���	cbz!�Ƨ�����⽏Ĵfr��9v���`\;�@]���쵔�)jV�      [   N  x���[jAE��U�?������,"+�V�HB�3"6�y�K�16��d33���TǙ�S=F���{�T}��d=��@�mr�ܠBDG�<�n��ݿ���{?��F�ao�7!� !ͭT�����\��|����"s�C���|��<��	=|@k�nl��j��(��1��Ze��������ೈ�2��m_��ώ��s��dw��m�Uf3��~'���6`]���o>����|2�ٗ�A�-g�4h�DWb$�<,�v�&���&�(�+�U�&;��e�+ǚ��$�ר1z[W���%�I�Q��5rOo�dr���^c�m��9oԘ�񸸝���y�9`�Z� ��Lϩ6��X�����CR�oӈ��F��F��r��f��`MT��W�[
L��6�A��G��co'��GZ9���Ѐ B��*�i��Z�XΦ�I�=<r�]nIDl�������br�oT�g��uK����ֵ���m��Y0�(�������M��"�w�3]��b|�:�)�>�����m�]�D򜠇�Q򐌻_���Yb7�f�!np�b������v\�O��T�~�/���      r   D  x�u�;n�0 �9>E�"��t����N���U��D��")_�(W�Ү�M@��s�z7�h�Ȓ�/z{��]h|�?�_��`��f��S����Q��Ww	H'�=�=�OF	�f����؉��=y�p�O2]��,�]w�������wπ�ak��K�dɰ��1�p4�ċ�5;�狧e7i-j/�4G?v�-��A�0,d���c�<��GN��^��֞�:�1�ً�kc�[����4��{����0j/Gy͏�1{���;�"4���1y��G\r�@b��`֨�X��Ʈ3��w��}�w�m�~)7�      t   �   x��ѽJ�@�z��e�{�ܙ;3O�جtA��7X�B�E���tk��4���
n i�7p��|�8��(+U���|�i����(��q���OgG��̄ G�#�&���^Afp7�)��VP��%�]
�8 J/�1��l� �I �$;F�N���6R�p�,��n�V���y\�������u�?��kz�����(Z�Q��p�tc�:�ܦ��lq(�P�3.�����-���8�EBp�فβ����      R      x��]ݏe�q6��C�DY]U]�q��Y�����a%ia�H<�66fQ�f!�2�&�`/�`�3sgg��/����̽u�9�w�\/#���]������A��q��;� XRXw)b�;��s��g�yn��փշ�>����>z퐼��{��2C���;��:��K�%��@
:�{�].`��ŋׯ����g�,��מ�������/�X=��z���?9����ݽ�X���������>Z}��?�Z,�ȸDZ���)�V��T��o~|�����X.�W����g������.&��sK���9L����
9X�?�-���OQsA-�|8��,.�� .\^�0u�]�'�|l��q��W�^Z�߫�o.���/?z�e�i�}R�/
tp���?�����=8���u�(6�	M9��/^Y�qG���^��8P``w�@ؼ��D��T����藔:�D�/�������AC�pC�m�X}������#u���˘�._���0�����<��'X^b�"E;�G�~�������L��p������[�o?-0i��R�,���-8�c�@�M,8��Č9�����y�-U���b��.w����Y��f��'�&�n����P��_�v��'�y��k/�x�����8oߩ�	8�%�~�����A4s��d�Es��Ls�28���LgO��S����LR�.�--�>4kߝ�`��('�*�g,��7-�0	t!����2������Vo����Oe�/m��ⶑo�Ìy�"^$L�2y��U���^h�����CV�`��,6�O�멞m*.�.t�.b�ƍ�C�! QB�)�m���O��[P�Ń�#������j�\�ܬ���[�ۭM5�
p*Ǒ[�k`m���ea/�a�J�fM�S��b,��J;�z*���?]�����?�����۟}��B�!����c���z�����}���x7�=�	r��1�����V�H��]t��ҍ���l��U�5��Ed�ՠC1��5d9��?�d�6@��Ĭ�2��`�����p��Dc�o9F�i�S0x�!�Xx��q��SO_��%?
v� �i.���t꿹V����~��wzȳ�>ZXY&�W@���9`������yk��-lM����ʕ�^{�ᷟ޹}��/�_5��HLȧ�~�ի�cM]�-���rs�}�9��0�PXj�$6���/�{��p(�E�+�m+����J�gr�k��� ��G�\c?�H�SU�X�x^�C�v�����7/��z�Oxq"q�H�k�އ?����oՈ��G*5 -�O\y���5"�������W��x�F:u*�LEN�4L؜٘Ľ�W9>|�����d��J�{ִ����?y����\{���<3�f�������na�~���4�i�޻�x����2�W��7;�i��&��]!7ӿ�,�th�*T�v�ѻ���@C@�ŋ�G���_�@je�N�?gO����x^��Ly��l�E ����V-`VWɣ�ٽ�,!@O�����+�� �gPtJθ)���xإ�3�e�6 �C] ����)V9ҋ=6+�3nk)��ڂ��|�~�7�S�oqM���G���H;^��/Dހ�N3��`�P���3ﺈ�9�q�ha���r
͹���&˴D����=�s�f��5�����695+]d�c����$�ͩ=!wX܀@#�|N�G�Jda�^qsp�g�-lT�Ap�Om����fq:H�C���a+F-B�a�*9�-�ܱۯ�n�(qQ�.�x}��=�蓣���@���A��W^[��WJ~M�r{��fJ��//��r��/����Z�!|	����`�'c���oݜ�ܛ�C@9Hb��s �nGͪe2@z�'���ܜ��me��[9�I.�f� q;"I���}�{YC�,�hG1�k���� �1���R�Ƅ9Y 	�RG,�E���Y�!��� i���5/��`�I�x_4�ݾ6pV�p�,���~u<���/nS��c�rs�v����IN.�L����y�D��� �/��[[H�9z�f�/�	Ps��N<�=-.ZHҫ2�4S�a��4G>��Ŧ��-[H�V Q���g���"q�8;�/�������<�cԗ���'�$�{r��/gv]��'Z��Ǌ�����X��D��Zv��M �>A��Z�:��G}�*X}�HQkHӞL8��>NIj�)7��]!���%y����w $�����r��n	⬣�בg,�w�H�3Ss�����!�sD�m����w�,RԜM�ܮPv]�Bz�q%6���{�� ��G�����Ņ-�I�3�Lެ���
�S)R�nζ��8�PR�y�3)h�c���"��ؓB!�P���]�ܞ�?cq��@V��=���E������BQ�$�k�;cm��@V���e�{�rdՉ/�_�cn��;cm��@u��fbv�'C@V�pІA���%�d��B�����z��6�I\B�8���#�..!��4d�������k��"!��/�v\��!$�n6�̽��7�	���R��]�����#��Xz��r��ִHZ�Tz]G�\s�ܮ���;�����3���֫�������{���ƩTk�^'޹o.�8cq�tR\L����c���Q��s���������Y�ȧ�_�Y��~6�xRX?�з{����I1qE*fMB��~�N
�I��k��_1p{� N�z�p|w鲫U�[��AY�HS�]�`�� ��BiH9���:(�A��P1�c$�IYlz�M�(��K�O�@{�凡��Ss�Lc����7o�,�y(k:;7WZ������A� �R�g��S��~��᧿.����y�36>�0G��/4�F�����j����J�
96_��Z��©��ԯ`Q��6��e����D� ;�ND�Ӈ0D��i��֦�qRԌh��WL/�wjB2HUuz�!nT߆�1H���Ƙ�%��e����Xt[PS�8R{M����QG���2v��Z��f�ޛK�,��9��˦ �۲&�tl�H;8!���E�D��䵊J����a��6yn��Eh��:n��e- uPcޡ�T[����E��5�.��p��H����W�O����O_��TEH��S
�i}3s�����j!E{����^{��5gik��=<�'��D;^�s����3e���Ќ)S���G-q";�(v��b�f:}�!f�8�ȦS���4����14�3�ۧ�i�ʰ�[�q�5�5O��dB:�9$�J�xC�=c�[
h�z�A�8�5j�C��&^A�'D���R�V)�=�1��f��5���Ɓ�"uE��qs]�ᇯ<�u�������kW�����/��@4���d��۶o�p�Js�Zn��c�M]���ؘ�1F����2��h.glق8��Jp�Ps�ۦ��� A䂚_�h�Z�0��z�$�9tΌ�&�e�CUK��lrh�4��a����1h�h}�b�ntnd�}$��T��6��1QB�yD�б�4Q���G�]0cؕ�f�0���C�?�]�8i8F�
����-$�ΎѤN�����"X>�M��X6�.m"#҃`ـsy+&aj#�L�A��#0���%�j�䜨�-������0�a�ʈ@�S�vL}�Ѕ��\Z���0ڥ���m3�v��bv�F�:;&��O	ipL�q�=}I͸�C55:�� ��D�qs�sî!Y(m���5z;Ɨ�Z���+-S�!�����1����r��9��7L-aTR�~�@�5X�+f�Y��rc3#�1�֥J@�����u���aϥh3z�B"$4���+[	c��6�9����)fi�f�� �[��ui8�!+38���ḨP����wH���e�)��Nf��jB6�h��t�l�Jڱ�w�ȣ��nc-�nӊl~	m��=�{}1�f�4EyL�j�� �  CsM�ؽ��O?=����ߞ�١'���(��x,CŅ�>i���<ǈJ˞q�l�U�y��@��d��������7oJ�|,f>Q�P��36�v�{2FՓDl��aqcj�J{�^o��gT�'�}{��w�.q��e� �|��M���!ב��g9&飥_��A?l�)�d�ǒesi�ސ��|$�[n/.��φ>�+V]ʮ�(z�~�W�>���;eζ=�^g�+T��B��:򉛓�-ϻr���u��6�O��!�w�y1��[aS��Y��C\��Z�)�`�'���܎�#���/���;�f�?��v@��i&�{K�8��14W�N�gK�T��\{�<E?X��ҔI�f�t���zQ��fa�d��d������$�<���a�w�V|%�A1}�9�?Z��d=w�b;����/���9�/�з�˾��$@����_.�Y����{���_5��ޔ;E�J/����1�v�Vz5��:�2�,��)��}��Y��Y�?[�^S$�k{m�����z��f�6�W�(
X�r�L�2���J8L8��JZF�E%2���k�bȡ���)���g�/��C?�}�g�qjm��5Ma�^��Z'ƅΛ��(�~u�f�/+ݘ�`�!����}+�ژI�;B{���l�[��	��d�{��oźwMe���kW���%j����Y�]���z(��,�!��Q�R���PFk7�=���T�lva�~��:�)����7E?�k��>><���������΢L�j�>��ɬ��;4A�ӧ����>���9ϑ����J�������裥�kJ
i�� �M�-�B�܋����yK��z���������wy)������O��'B��&LҏC�z�/~i�`������8Rhe�~���~�Kz�?���Y�Z��EqM�o�胥��)��;B'�o�o͈�N�c��d��d�sD�~H_|����x���-��~<�"�,}-��.Kh����-��wi�˞.J?��{C���,v�~��Q31����V���k!�6������׮-����}+�Z� �}������ߒ��D�`�+!B+����^4��V~Q�X���#�]Z��q�%*�#�"��y�,���3B+���Қ�F�g�+�X�r#�>���&蓕_�^ !6';&�[���P ��E�)�V~��)t�|�Y��򋥚&t�eS���ֶM���A'�[�e}��s�oNRߔޓ�h�ln�xw����{��<iI       j      x������ � �      b   �  x����jA�ϽO1w���?��x�!�@�$�=i��bT�{�	*�i2��-R=�dwfV"+[{�.�����5���Q	P�TiX ��w\�%�}�4_�]�8h�����|��b>{f �FL�0"@J�����oZR�XK���Y��Ƙ����<z�[��<�oG���5'/M��S���P�<00h�T!�HY�$ì�\����Ϛ/z$H,>Ĵ� G��Ղ�\ys�>h�@x)[���!`���BS�E���6���c��HXT7�(�QB�1߱�st��I���P " H.��3
D%+�����|t�AGa��S\�[��~k���W��4?9��������n�1�i�U�a)�
�lJV{�{��6٭GתS`Y�ھ��Z,��:�2tN�8�D �VŒ���D�C"�^s��*�Ĩ��/�v����<J��b��n���B����c
�J�J�H��y==��}���4n?�(_d���Ze%��d�N&�k�0      I   +  x��Z�nG]K_�e�q*uo���.���/�	c��� �NB��C���H$&�=�+�,*�ϰ��j~!��i6�%>$��5 Q���ԩ[�Q�0��'?�ƣ��2ͥT��v&Y\l����ݍ�ۯ��k��j�m=�_�V[�����<H�X�w:o�W��Ϛ�����a�Uke��C��܀3tuf{���4^*`��C̰�?�+O��l,�I���^{���i�^|Rm��[On&W��_� ���v��R�v�D%3�U��s",r²ܣ�J3P �W������i�\�"��B*m�%<Y���W�)��I�9�S%<O�����)<j&��f�!�Ҕ %x@&�E0S���s�ǅ�O�[��r7�#����3I�!F���b�zшD2��ҹ����c���ܺ�}��d�Q��{�w�-%[�m�Wۛ��J�k�Vg�_8�տϞ�u�i��G=�Ѵ�7���������,��}h�$N\��8�V�6���ri(cp/��2��`֧M���V�bsk!Y���ߛփ��ҋ����x�h��a$ٸ�j6�hׯ�+��j��T��R-3 ���z���*'zCl-l$?�u�ﮖ��/����,�)��l���N
��n�4��z��[Z� 
f:(�;m�b��TDTG�1S3C�Ԍ�����uꦎJ�rF��O�z	��6�bt���+�3�l\�a��ן�n/�?�H᱘�7�j/��[ �ϒ-�����B�*@?�Y��{3�}��ʳd�ZN�b����]�_����c�Ί*q*�)�X���<��,]���t�#ѤG�����p��ԝ���E�=�$�i��ɳ�0yY O��T������H?z}��U���;�Zm��2M^�ɛ�a�+k�����z?�M����4c�)���V��&/�*0aQL����w{�WLq
�������0{�E_�7鶛�UN��=����D/\�!fX63qkQ���
Q6�p�LID#��ɒ�J�JJ@3��`��_%M.X���\1@Ie�3S6��0a���8ڂ���[g� 5�,<�kB���*9Br}���zҭ�Ϳ��Ka�n�F.�� B�����%'��d�%�-�O��ehCK�)e5NYp�%hZkT��j�
��(#�,s �\'-��*�}j���a6!hU��1HI��ེ	A�����E
ӆ6%h/��ཿ�ȱkn]i�$w/7/�G1�||���gϝ���oN�=}��O���>7rNQ��Ұ��~�ݩH��I�{iJ((OYW���,HP��C	C�PI�N�."���6mD�a�=N[���n/�a�����v��y�����;��S�z���s��ϒGO��>Q��"���,~׶;�6�k��ڽ�rf.=�1^!��۴k��ݰ]���U��
{� �3o���U�T�9bm����oU�kJ�qJD	2n���ƹq!;�k��z����������)�1㲵^*�ț�����8ўD��:�C~�`4�)=�y��W��w*���AW�Hp(�`�7V���]k8CJ�'�Z5����i�W+�7]�0��S3��c�T�S͞K:�̶���㥰��7'�g�
����禱��Eԇ�)�taw5�?ވ�a���#��m����Գ�5����i�f�T��X��r{�l��h���&s�r�S1✧7S:/$�:=\��:�-��pK�W�çK�"6񥟼J@rS��,��#�`v�t�"f�y�N��ib��4{�'��Q��-����F0�fZ��X�rdU1T�;:M��,����9&��v:/�����S�����6;;�7��G�      G   �   x���1
�0��9�������i�	z�
���Z7)��i��l]��_~��#�0�p��%HB)>����̻祻��¬��O�T�i�U��y��}Ä�t��j)s*�C�(
�3�:L�RM��t��S)9��Z�{+K)�X�TAi��R�s�T$����u�      n   k  x���AJ1E�ݧ��L��T%U9�'��'q�q!���#�efF�aFӃ������T>q" X�_�P���ՙD� ������x�=�<~>_��w���4]M�����a��(&O�P�p���m1R�q|l����������߼$���B�1� bq`����d8�T󘈝�o����0��a����H�:� �!m�
P/k�K~MpD�(��8��`랡�É'H���8�:���_��ĳ��qe���P*���a��:���� ���Pu�CC����]�����D5��Gl;T�<�ΰ�Z��(��p��yk;4
����C�\'E�u(��v��]��7����<      p      x������ � �      l      x���ݪ-Gǯ�y
_@��ڋ����^�A��Q	1�Q����<�O�}r�Ξ��5Uk�׮��I�bC��������=���[߿D{���%�����ُ/������?ߵ���������G����:��ń��3�"~����/޹�y������7_��𻯾��7/��s��jL�5�q[@�<|����/���͗��,����]g�]h�Җ�+M��㿿�����{����_o>����?��Z�I֬$�H2�%���'�0���~���_�����c��ٚd�n���<�s3[3�Ӝ�9b�3����G�~���_��@��٭�n�#�� �Zuxħz����/����E8C�:��:�\J�˿XZ���:GI���Zo��?��}���@�34���_��w���+�*�X��ҝ�g���Z*w`U�鼒�V��㥚y�9B��4v�'��b�X�%���+�R>{"v��Evv;!�R{����3�eO���W�Wx�pV��٘1f(l�i���>_�d=��-/�כ�V��t���a�P���t��nt

$�Q� �;�*�(A�X���KV0/�,-���	|�F%����Gl5+l���h��W�1B�����
0B�X���sA��6��'f�o�ӥ�`��a.�s�$�Z�-�0@O���/���H`�+?��	���#c@�߼|�����
��}�fd-DU���|��U�PĜ��R�0� �!x�6.:�ĺA�1�c
u�Wo�����P0`��A�%?O�$Z ;�-E��]wp�j�#z�e�cb �'͌�]��O1a�$�3fL�mo'��즶d�,Gj�a+��8����~R�$�s�H=^�H<o��k�	��=iS����=|�� �c���̑��,���`��9�P&6|���U<r7@� 5�u�G�?҄9�G����ٮ�fR�5�w+����׽��Pڍ���[�\*��>�y�*6�� -Ov�L��L1g�n��v7w�����Y_0@������9R�`�����g���Ǖ�)��!���&����FJ�pbݮ2��G�P0`RHh������`<�/G�s�F���[�<ELV)P��!����߸�p��v�$���}�}ܯ��l��1W�d0M�;N=�=�?qn��x��<L�	��G�s�]8s�"�y�9a:'�� X����sn����0�59[�6e��eX����H�P�)��>����	|6D��j��ĺ��ﴂ���B�O�����JUI�*���⯫��a<3��3Ě�K��;C��>b3(a�ؠ�9b�
����6s!�$��*Z%�m��U����2mG�7��{Jί�s�<�E�� �|QE6T�6�X�Ecr-$�du*��L=����a��b�`��aƔ���W?g��L�x�(LĀIk����v��8��uG� g~� �H� �$��@���ɒ'�	!�'�6��y�����|��.����/�
��
pFHGxX�]-�Zׂ�V� X�]�/�3�s۾��-N���X�ʵ=nI�6�>��{ܒ `�X܈9b���	�W38�8WC��V)I����k�_� K@���*�w��t����$�Ts�֩����k�����#a��Ԍ9�6��9�n=�fn�x�����<�W��#���$�=�S��[ڤ� �[��"���ŉ҆�4�X�L@�0��k��=�^7i��S�Y��.�g��[(�=;^7j�<���('�q�����^m�o�b�DCLK����!����>�J[�Q���I%�gO*�MZkL	Q�v��<���hr�y$EZ {�����1m,O��uW��b:w�ʇe�4f����8����"���=(��G�fsE2�5F�����+�K5�+6�Tr8T9��p�\�C��z���1^�W��^0��H�LJ�7�2-���-.�v;��W�\*u���T˺�W���tiÂC�g(���l�p��^�=�S����W6�0@o�����O�`�Xp���h5�z+꫷V��HU�.�$0�'A����Y���b�x�׹���U�J2Eh��B��6��gK��W�n�%�93l�k4��1M⨦�ʎj���a�Ô0@�1L�'89�)Nܔ��^ \�e/�J��$������N�81�A�[��Zω�w��4��[�7|g��߅�1�5��9s�z^Yū�V@�M�KNs�--���:Q���p:zCx�`���J���佡O6��0G�;6�2�sF�( 큥b�P�ڏGLhr욮6U��׹HE�s�mަ���+)R���ܤK� �#��"���7j���[+�>�s̟�!��&c	H�.�1;9�ݺ��`��L��1f�5��#��#��yl�<��n�e������!p�v6Ë���2�}b��N�	^������c��2=83a����E١��3`�eu��]~��
���b	����KuW��1^���Ƶ�˜{�#�����]���Vۀ��E̘4���Ƕ�f��.�+0im�N���"y�9�Iwϱ��7n�G��xθ� ����?'�
� ����"�O5Q�����s^�t����G˖���R��! �R�������jpp�p�c�X��9b�z@~R��!���0�3d7@� �^s�z�2@�t�jp4� �R�!,E�c�;����ۥ�>��ɥR���|d1b��⩉�0��xj���挟4c����U���֮��R��>x5�KN�6`Q�?�� t��\���5�z%L��_�*3�S� eO.0R�&�7l��6��B�R���-`���r��_����cz�9��Jp��#�]|����>g.?� `����Z�b0@�],��-�ﴦ�. O@��Y��@@�/�U�H�O�|�凎 9`�12f��.�#ֺb��}���z��b����0G���f��x��� ��n�H�����%J��x�<���'U(I���o�{��J��`:���� �9�8�=�l<N,|x������6�9$7@� ��s�z�v���������yEXCR��1^r�"A��}H|3�c�rC�n�JV"k!�.8"�U(Q�rC�n��Ϡ#�!��=ܽէ��:�dv�p�̄��е��a���|����>�F[��̌>��2|Gf��cO._7iU��HpF{�P�](�Gf�����f�̌?='�����>���_��k'�|zF�=��f��y��$�+��w!%J��{��C���x+����2{�+���lw��ث��J�m0z5��z�Gj��t�36���Ҷ�^��/ ���ഏ ���H�o����S0G,8��I=��b� ��p�HR���+�8��,C�]�Fo��?Fo�H�b��ņX�K���U3��_1^�|Vf�S��D��q�#����V�*p
��U��ໜ%{oN���J����qU?�!����@{�P6��d��s��AV���}b��?���0@�c��'��Ќ��ݭB�"�]�\h�6Ct&��j�vLϯ4WK�sf§c�]|8�§c�z��D_�Qi��]�{c�ٴߕ���dZ����X3������[+T)B/�� ߛ�O���fYm]�]h��n#��&g��$XͲ/_�t�%j�����	s��f�X�-,\��Rm�W��Ly�c0^(n1s�zkU��T�I�"��cZ� V-R�X7X�椿ٔI��-��-�+���T)�ֹo-���0��K��l=f��R5<�g�wjn��z�J�0^�H3�)x�k�6�*E(b���/;��`�h ����\6ܹ�f�K/�h	�"�I��K�����'_�o=���'���������{��(�R�xG��{�y���ʾ	m�� ��B�	�'�ɘ6V|}�ݠ[�+(���`��4���0�R�8�l�
9��+�<-���m��B�vi�'�49kCbO�m:v4��1@-����Y�7{3+�Ս'�j��u��q<|�>:Ћ�7��x��p��^}�:ݪ�r�;�ql���)B/6�
 �   (� ��}��`�y���;�����$�G���7�#\�Rp��x�;�����7|�x��	s�m�;��&����fP� �����/�+̑��b��"��$8��/�jxZ�`Ceos'XѸ'a{~�q#'4��DWS�*�P=r��9��\��b{$)�G�|�N��:��aiϽ��(�0�L�=c�|��#<[-���s��d� W���j0m��@<Tf�����|}��XC�����ŋ���      K   X   x��̻�0�ڙ�>2zNBl3����#��'�PB#q�&�;�������������K��MGa��ݎ����U�p�!��M?      P   �   x�3�|�zƳ+_,\�b�NCSN��Z.#Η�^�[�b�NN#��1H��u9-�"F �';�>���i1�|�w�4NC�.3�DNC 2�9Ӏ ε�L,NIK" D0qKΧK7?ݾ�i���S惜 2*Ə�h	Y�y���;�=[��|�>�
#Χk'@ŀ
�j��1z\\\ �T�      ^   �   x����JA�뽧�^n������>�O��@�N#	B�-��� ��!6"Dėɭ>�'9T,S��}�S�`d��5���a��H��M*ȭS&��t�o��Y��}��|>��q�n�G㣪60����U�IUW�A;[���x�/��n��vZO��N\�-&os.Kҍ,� �����.����>��S,�P������CJ.�*m�k�7kM�W�(k���4�����m;�wm�$_}^��      d   J   x�? ��1	2020-02-20 11:17:34.701019+08	基本工资账套	{2,41}
\.


T��      h   �   x�3�4202�50�54R04�26�22г07701�6�  ���lچ�3�=[��Ӏ�@�����D��������D�������<|�ƺ�f
�V&FV�z�F��`�g�k���m�өm/'o{��bg� ��2�      U   K  x�eйn�@�z�t������*��p��1� C 9<}0�($�H���o~��@����J��%I�$��e��
����4e�n�`LJ�`��� ���gAdG����Zu�s�L3n%�c�:��K��d������Ͱ[�ö1�O�hv1{gLX�}s�r\���#)	.��ȸ�pP�
������:�l��G'�}��j؜5�VT
�m�ͮ�x����8�pz�q'	td����Ȣj�m�\�(mK�[J"s���47H��ͿS�Tw�ZT�������}�L�??����&�-���<1_;N
�[�oܲ�/c�t�      Y      x�3��4����� �Z      \      x�3�4�211z\\\ U�      V      x������ � �      D      x�3�4�2�4�2������� �      L      x������ � �      S      x�3�4�2b 6�=... &�     